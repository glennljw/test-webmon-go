package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mini-webmon/internal/api"
	"mini-webmon/internal/record"
	"mini-webmon/internal/retention"
	"mini-webmon/internal/store"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const sixMonths = 180 * 24 * time.Hour

func seed(s store.Store) {
	for _, line := range []string{
		"2026-10-01T10:00:00Z|https://a.gov.sg|up",
		`{"target":"https://a.gov.sg","result":"down","ts":"2026-10-02T10:00:00Z"}`,
		"garbage",
	} {
		rec, err := record.ParseLine(line)
		if err != nil {
			fmt.Println("Error, skipping line: ", line, "e:", err)
			continue
		} else {
			s.Add(rec)
		}
	}
}

func main() {
	s := store.NewMemory()
	seed(s)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	job := retention.Job{Store: s, Window: sixMonths, Interval: 5 * time.Second}
	job.RunOnce()
	job.Start(ctx)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           api.NewRouter(api.NewHandler(s)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("Listening", "Addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Shutdown Failed", "err", err)
	}

	slog.Info("Server stopped")
}
