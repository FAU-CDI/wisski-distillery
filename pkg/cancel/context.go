package cancel

import (
	"context"
)

// WithContext executes f and returns the returns the return value and nil.
//
// If the context is closed before f returns, invokes cancel and returns f(), ctx.Err().
//
// In general, WithContext always waits for f() to return even if cancel was called.
// As a special case if a closed context is passed, f is not invoked.
//
// allowcancel must be called by f exactly once, as soon as the cancel function may be invoked.
func WithContext[T any](ctx context.Context, f func(allowcancel func()) T, cancel func()) (t T, err error) {
	t, _, err = WithContext2(ctx, func(start func()) (T, struct{}) {
		return f(start), struct{}{}
	}, cancel)
	return
}

// WithContext2 is exactly like WithContext, but takes a function returning two parameters.
func WithContext2[T1, T2 any](ctx context.Context, f func(start func()) (T1, T2), cancel func()) (t1 T1, t2 T2, err error) {
	// context is already closed, don't even try invoking it.
	if err := ctx.Err(); err != nil {
		return t1, t2, err
	}

	cancelable := make(chan struct{}, 1)

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer close(cancelable)

		t1, t2 = f(func() {
			cancelable <- struct{}{}
		})
	}()

	select {
	case <-done:
		// the function has exited regularly
		// nothing to be done
	case <-ctx.Done():

		// context was cancelled
		<-cancelable
		cancel()

		// still wait for it to be done!
		<-done
		err = ctx.Err()
	}
	return
}
