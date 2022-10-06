package timex

import (
	"context"
	"time"
)

// SetInterval invokes f with the current time and then spawns a new goroutine that runs f every d, until context is closed.
func SetInterval(ctx context.Context, d time.Duration, f func(t time.Time)) {
	f(time.Now())

	go func() {
		t := time.NewTicker(d)
		defer t.Stop()

		for {
			select {
			case tick := <-t.C:
				f(tick)
			case <-ctx.Done():
				return
			}
		}
	}()
}
