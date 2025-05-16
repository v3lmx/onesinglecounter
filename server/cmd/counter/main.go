package main

import (
	"flag"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/log"
	"github.com/v3lmx/counter/internal/api"
	"github.com/v3lmx/counter/internal/backup"
	"github.com/v3lmx/counter/internal/core"
)

func checkCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "Same-Origin")
		next.ServeHTTP(w, r)
	})
}

var count atomic.Uint64
var best = core.CurrentBest{}

var port string
var backupCurrentPath string
var backupBestPath string

func main() {
	log.SetDefault(log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.DebugLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
	}))

	flag.StringVar(&port, "port", "9000", "Port to expose the api")
	flag.StringVar(&backupCurrentPath, "backupCurrentPath", "./current.bak", "File backup of the current value")
	flag.StringVar(&backupBestPath, "backupBestPath", "./best.bak", "File backup of the best values")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	commands := make(chan core.Command)

	var m1, m2 sync.Mutex
	tickClock := core.NewCond(&m1)
	bestClock := core.NewCond(&m2)

	backup, err := backup.NewFileBackup(backupCurrentPath, backupBestPath)
	if err != nil {
		log.Fatalf("Could not create backup: %v", err)
	}

	backupCurrent, backupBest, err := backup.Recover()
	if err != nil {
		log.Fatalf("Could not recover from backup: %v", err)
	}

	count.Store(backupCurrent)

	best.Lock()
	best.Best = backupBest
	best.Unlock()

	go core.Game(commands, &count, &tickClock)
	go core.BestLoop(&count, &best, &tickClock, &bestClock, backup)

	api.HandleConnect(mux, commands, &count, &best, &tickClock, &bestClock)

	log.Info("starting server on port 9000")
	log.Error(http.ListenAndServe(":9000", checkCORS(mux)))
}
