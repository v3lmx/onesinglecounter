package bootstrap

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VictoriaMetrics/metrics"

	"github.com/v3lmx/counter/internal/api"
	"github.com/v3lmx/counter/internal/backup"
	"github.com/v3lmx/counter/internal/core"
)

type Config struct {
	Port              string
	MetricsPort       string
	BackupCurrentPath string
	BackupBestPath    string
	CounterTick       int
	BestTick          int
}

type App struct {
	Count           *atomic.Uint64
	Best            *core.CurrentBest
	Commands        chan core.Command
	TickBroadcast   core.Cond
	BestBroadcast   core.Cond
	Backup          core.Backup
	CounterTickTime time.Duration
	BestTickTime    time.Duration
	Mux             *http.ServeMux
}

func ParseFlags() (*Config, error) {
	config := &Config{}

	flag.StringVar(&config.Port, "port", "8000", "Port to expose the api")
	flag.StringVar(&config.MetricsPort, "metricsPort", "8001", "Port to expose the metrics")
	flag.StringVar(&config.BackupCurrentPath, "backupCurrentPath", "./current.bak", "File backup of the current value")
	flag.StringVar(&config.BackupBestPath, "backupBestPath", "./best.bak", "File backup of the best values")
	flag.IntVar(&config.CounterTick, "counterTick", 5, "Counter tick time in milliseconds")
	flag.IntVar(&config.BestTick, "bestTick", 200, "Best tick time in milliseconds")
	flag.Parse()

	// Validate configuration
	if config.CounterTick < 1 {
		return nil, fmt.Errorf("counterTick cannot be lower than 1")
	}
	if config.BestTick < 10 {
		return nil, fmt.Errorf("bestTick cannot be lower than 10")
	}

	return config, nil
}

func Initialize(config *Config) (*App, error) {
	app := &App{
		Count:           &atomic.Uint64{},
		Best:            &core.CurrentBest{},
		Commands:        make(chan core.Command),
		CounterTickTime: time.Duration(config.CounterTick) * time.Millisecond,
		BestTickTime:    time.Duration(config.BestTick) * time.Millisecond,
	}

	var m1, m2 sync.Mutex
	app.TickBroadcast = core.NewCond(&m1)
	app.BestBroadcast = core.NewCond(&m2)

	backup, err := backup.NewFileBackup(config.BackupCurrentPath, config.BackupBestPath)
	if err != nil {
		return nil, fmt.Errorf("could not create backup: %w", err)
	}
	app.Backup = backup

	backupCurrent, backupBest, err := backup.Recover()
	if err != nil {
		return nil, fmt.Errorf("could not recover from backup: %w", err)
	}

	app.Count.Store(backupCurrent)
	app.Best.Lock()
	app.Best.Best = backupBest
	app.Best.Unlock()

	mux := http.NewServeMux()
	mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	app.SetupRoutes(mux)

	go app.SetupMetrics(config)

	app.StartGame()

	return app, nil
}

func (app *App) SetupMetrics(config *Config) {
	testMetric := metrics.NewCounter(`osc_test_metric{label="test_label"}`)
	testMetric.Set(10)
	testMetric.Inc()
	testMetric.Inc()
	testMetric.Inc()
	testMetric.Inc()
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	slog.Error(http.ListenAndServe(":"+config.MetricsPort, nil).Error())
}

func (app *App) SetupRoutes(mux *http.ServeMux) {
	app.Mux = mux
	api.HandleConnect(mux, app.Commands, app.Count, app.Best, &app.TickBroadcast, &app.BestBroadcast)
}

func (app *App) StartGame() {
	go core.Game(app.Commands, app.Count, &app.TickBroadcast, app.CounterTickTime)
	go core.BestLoop(app.Count, app.Best, &app.TickBroadcast, &app.BestBroadcast, app.BestTickTime, app.Backup)
}
