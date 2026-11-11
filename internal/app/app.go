package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkrav3ts/thefinals-autoclickers/internal/keyboard"
)

// Run calls fn and returns when fn returns or the process is interrupted.
// Windows timer resolution is raised for the duration so short sleeps match
// the delays the tools ask for.
func Run(fn func(ctx context.Context)) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	keyboard.WithTimerResolution(func() {
		fn(ctx)
	})
	if ctx.Err() != nil {
		fmt.Println("\nExiting...")
	}
}

// Sleep waits for d or until ctx is cancelled. It reports whether the full
// delay elapsed.
func Sleep(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return ctx.Err() == nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
