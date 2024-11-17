package main

import (
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/v3lmx/counter/internal/api"
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
	logger := log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.DebugLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
	})

	mux := http.NewServeMux()

	api.HandleConnect(mux, logger)

	logger.Info("starting server on port 8000")
	logger.Fatal(http.ListenAndServe(":8000", checkCORS(mux)))
}
