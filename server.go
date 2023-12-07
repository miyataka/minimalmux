package main

import (
	"net/http"
	"os"

	"log/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	srv := &http.Server{
		Addr: ":8080",
	}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server ListenAndServe error: ", err)
	}
}
