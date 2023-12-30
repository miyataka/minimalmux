package minimalmux

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"log/slog"
)

func RunServer() {
	config := newConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: config.LogLevel(),
		}),
	)

	mux := NewServeMux()
	mux.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		healthcheckHandler(w, r)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: mux,
	}

	if err := ListenAndServeWithGracefulShutdown(srv, 5*time.Second); err != nil {
		logger.Error("server ListenAndServeWithGracefulShutdown error: ", err)
	} else {
		logger.Info("server shutdown gracefully")
	}
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
