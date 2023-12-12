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

	mux := NewServeMux()
	mux.routes = []Route{
		{Path: "/healthcheck", HandlerFunc: healthcheckHandler, Method: http.MethodGet},
	}
	mux.setupRouteTree()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server ListenAndServe error: ", err)
	}
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
