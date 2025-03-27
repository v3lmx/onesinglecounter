package main

import (
	"net/http"
	"os"
	"sync"

	// "syscall"

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

var count = core.CountM{Count: 0}

func main() {
	log.SetDefault(log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.ErrorLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
	}))

	mux := http.NewServeMux()

	commands := make(chan core.Command)

	var m sync.Mutex
	cond := core.NewCond(&m)

	go core.Game(commands, &count, &cond)
	// go core.Game(events, clients, count, requestBest, responseBest, cronBest)
	// go core.Best(count, requestBest, responseBest, cronBest)

	api.HandleConnect(mux, commands, &count, &cond)

	log.Info("starting server on port 10001")
	log.Fatal(http.ListenAndServe(":10001", checkCORS(mux)))
}
