package main

import (
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/log"
	"github.com/v3lmx/counter/internal/api"
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

func main() {
	log.SetDefault(log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.DebugLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
	}))

	mux := http.NewServeMux()
	mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	commands := make(chan core.Command)

	var m1, m2 sync.Mutex
	tickClock := core.NewCond(&m1)
	bestClock := core.NewCond(&m2)

	go core.Game(commands, &count, &tickClock)
	go core.Best(&count, &best, &tickClock, &bestClock)

	api.HandleConnect(mux, commands, &count, &best, &tickClock, &bestClock)

	log.Info("starting server on port 9000")
	log.Error(http.ListenAndServe(":9000", checkCORS(mux)))
}
