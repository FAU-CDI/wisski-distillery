package wait

import (
	"context"
	"time"
)

// Wait repeatedly invokes f, until it returns true or the context is closed.
// The invocation interval is determined by interval.
func Wait(f func() bool, interval time.Duration, context context.Context) error {
	// create a new timer
	timer := time.NewTimer(interval)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()

	for {
		if f() {
			return nil
		}

		// reset the timer, and wait for it again!
		timer.Reset(interval)
		select {
		case <-timer.C:
		case <-context.Done():
			return context.Err()
		}
	}
}
