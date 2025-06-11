package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/v3lmx/counter/internal/bootstrap"
	"github.com/v3lmx/counter/internal/observability"
)

func checkCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "Same-Origin")
		next.ServeHTTP(w, r)
	})
}

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))

	config, err := bootstrap.ParseFlags()
	if err != nil {
		slog.Error("Failed to parse flags: ", "err", err)
		os.Exit(1)
	}

	app, err := bootstrap.Initialize(config)
	if err != nil {
		slog.Error("Failed to initialize application: ", "err", err)
		os.Exit(1)
	}

	slog.Info("starting server", "port", config.Port)
	slog.Error(http.ListenAndServe(":"+config.Port, observability.TestCounterMiddleware(checkCORS(app.Mux))).Error())
}
