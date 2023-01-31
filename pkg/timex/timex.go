// Package timex provides Interval and Wait
package timex

import (
	"context"
	"sync"
	"time"
)

var tPool = sync.Pool{
	New: func() any {
		timer := time.NewTimer(time.Second)
		StopTimer(timer)
		return timer
	},
}

// NewTimer returns an unusued timer from an internal timer pool.
// The timer is guaranteed to be stopped; meaning a call to timer.Reset() should be made before using it.
func NewTimer() *time.Timer {
	return tPool.Get().(*time.Timer)
}

// StopTimer stops the given timer and drains the underlying channel.
// This prevents it from firing, until a call to Reset() is made.
//
// If the timer is not running, StopTimer does nothing.
func StopTimer(t *time.Timer) {
	t.Stop()

	// try to stop
	select {
	case <-t.C:
	default:
	}
}

// ReleaseTimer stops t and returns it to the pool of timers.
func ReleaseTimer(t *time.Timer) {
	StopTimer(t)
	tPool.Put(t)
}

// TickContext is like [time.Tick], but closes the returned channel once the context closes.
// As such it can be recovered by the garbage collector; see [time.TickContext].
//
// Unlike [time.Tick], immediatly send the current time on the given channel.
func TickContext(c context.Context, d time.Duration) <-chan time.Time {
	if d < 0 {
		return nil
	}

	ticker := make(chan time.Time, 1)
	ticker <- time.Now()
	go func() {
		defer close(ticker)

		timer := NewTimer()
		defer ReleaseTimer(timer)

		for {
			timer.Reset(d)

			select {
			case tick := <-timer.C:
				ticker <- tick
			case <-c.Done():
				return
			}
		}
	}()
	return ticker
}

// TickUntilFunc invokes f every d until either context is closed, or f returns true.
// f is invoked once immediatly when the timer starts.
//
// TickUntilFunc blocks until f is no longer invoked.
//
// Returns the error of the context (if any).
func TickUntilFunc(f func(t time.Time) bool, c context.Context, d time.Duration) error {
	context, cancel := context.WithCancel(c)
	defer cancel()

	for t := range TickContext(context, d) {
		if f(t) {
			break
		}
	}
	return c.Err()
}
