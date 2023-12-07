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

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", config.Port),
	}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server ListenAndServe error: ", err)
	}
}
