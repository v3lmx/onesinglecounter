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
		// origin := r.Header.Get("Origin")
		// if slices.Contains(originAllowlist, origin) {
		// 	w.Header().Set("Access-Control-Allow-Origin", origin)
		// 	w.Header().Add("Vary", "Origin")
		// }
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// var count = core.CountM{Count: 0}
var count atomic.Uint64
var best = core.CurrentBest{}

func main() {
	log.SetDefault(log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.ErrorLevel,
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
	// go core.Game(events, clients, count, requestBest, responseBest, cronBest)
	go core.Best(&count, &best, &tickClock, &bestClock)

	api.HandleConnect(mux, commands, &count, &best, &tickClock, &bestClock)

	log.Info("starting server on port 10001")
	log.Error(http.ListenAndServe(":10001", checkCORS(mux)))
}
