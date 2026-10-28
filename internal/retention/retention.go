package retention

import (
	"context"
	"log/slog"
	"time"
)

type Purger interface {
	PurgeBefore(cutoff time.Time) int
}

type Job struct {
	Store    Purger
	Window   time.Duration
	Interval time.Duration
}

func (j Job) RunOnce() int {
	cutoff := time.Now().Add(-j.Window)
	removed := j.Store.PurgeBefore(cutoff)
	slog.Info("Retention Purge", "removed", removed)
	return removed
}

func (j Job) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(j.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				j.RunOnce()
			case <-ctx.Done():
				slog.Info("Retention Job Stopped")
				return
			}
		}
	}()
}
