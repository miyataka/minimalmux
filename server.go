package api

import (
	"fmt"
	"net/http"
	"os"

	"log/slog"
)

func RunServer() {
	config := newConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: config.LogLevel(),
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server ListenAndServe error: ", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
