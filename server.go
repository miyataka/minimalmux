package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// prepare graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, os.Interrupt, os.Kill,
	)
	defer stop()
	ctx, cancelCauseFunc := context.WithCancelCause(ctx)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("server ListenAndServe error: ", err)
			cancelCauseFunc(err)
		}
	}()

	<-ctx.Done() // wait signal

	// shutdown
	ctx, cancelFunc := context.WithTimeout(
		context.Background(), // shutdown context
		5*time.Second,        // TODO configurable
	)
	defer cancelFunc()

	// shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error: ", err)
	} else {
		logger.Info("server shutdown success")
	}

	// check timeout occurred or not
	if err := context.Cause(ctx); err != nil && errors.Is(err, context.DeadlineExceeded) {
		logger.Error("server shutdown timeout error", err)
	}
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
