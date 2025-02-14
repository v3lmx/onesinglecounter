package main

import (
	"net/http"
	"os"

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

func main() {
	log.SetDefault(log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.DebugLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
	}))

	mux := http.NewServeMux()

	events := make(chan core.Event)
	clients := make(chan core.Client)
	count := make(chan uint64)
	requestBest := make(chan struct{})
	responseBest := make(chan core.CurrentBest)
	cronBest := make(chan core.CurrentBest)

	go core.Game(events, clients, count, requestBest, responseBest, cronBest)
	go core.Best(count, requestBest, responseBest, cronBest)

	api.HandleConnect(mux, events, clients)

	log.Info("starting server on port 10001")
	log.Fatal(http.ListenAndServe(":10001", checkCORS(mux)))
}
