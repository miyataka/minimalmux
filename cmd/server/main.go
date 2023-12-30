package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/miyataka/minimalmux"
)

func main() {
	RunServer()
}

func RunServer() {
	config := newConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: config.LogLevel(),
		}),
	)

	mux := minimalmux.NewServeMux()
	mux.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		minimalmux.HealthcheckHandler(w, r)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: mux,
	}

	if err := minimalmux.ListenAndServeWithGracefulShutdown(srv, minimalmux.GracefulOpts{TimeoutDuration: 5 * time.Second}); err != nil {
		logger.Error("server ListenAndServeWithGracefulShutdown error: ", err)
	} else {
		logger.Info("server shutdown gracefully")
	}
}
