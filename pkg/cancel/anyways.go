package cancel

import (
	"context"
	"time"
)

// Anyways behaves like context.WithTimeout, except that if the Done() channel of ctx is closed before Anyways is called, the returned context's Done() channel is only closed after timeout.
func Anyways(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	// context is not yet cancelled => return as-is
	if err := ctx.Err(); err == nil {
		return context.WithTimeout(ctx, timeout)
	}

	// create a new anyways
	any := &anyways{
		done:     make(chan struct{}),
		parent:   ctx,
		deadline: time.Now().Add(timeout),
	}

	// start waiting for the timer (or the cancel to be called)
	finish := make(chan struct{})
	go func() {
		t := time.NewTimer(timeout)
		defer t.Stop()

		defer close(any.done)

		select {
		case <-t.C:
		case <-finish:
		}
	}()

	return any, func() {
		close(finish)
	}

}

type anyways struct {
	done chan struct{}

	parent   context.Context
	deadline time.Time
}

func (a anyways) Deadline() (deadline time.Time, ok bool) {
	return a.deadline, true
}

func (a anyways) Done() <-chan struct{} {
	return a.done
}
func (a anyways) Err() error {
	select {
	case <-a.done:
		return context.DeadlineExceeded
	default:
		return nil
	}
}

func (a anyways) Value(key any) any {
	return a.parent.Done()
}
