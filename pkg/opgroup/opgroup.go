// Package opgroup provides OpGroup
package opgroup

import "sync"

// OpGroup represents an operation group that can send messages to the waiting goroutine.
// The zero value is not ready for use, use [NewOpGroup] instead.
type OpGroup[M any] struct {
	wg sync.WaitGroup
	c  chan M
}

// NewOpGroup creates a new OpGroup.
//
// The internal buffer size for messages will be expectedSize.
// If unsure about buffer size, 0 is a valid choice.
func NewOpGroup[M any](expectedSize int) *OpGroup[M] {
	return &OpGroup[M]{
		c: make(chan M, expectedSize),
	}
}

// Go schedules a new operation (implemented by worker) to run in a separate goroutine.
// worker is passed a send-only reference to the message channel which it can uszxe to send messages to.
func (op *OpGroup[M]) Go(worker func(c chan<- M)) {
	op.wg.Add(1)
	go func() {
		defer op.wg.Done()
		worker(op.c)
	}()
}

// GoErr is like Go, except that once the operation is finished, it writes the returned error into dest.
func (op *OpGroup[M]) GoErr(worker func(c chan<- M) error, dest *error) {
	op.Go(func(c chan<- M) {
		*dest = worker(c)
	})
}

// Wait returns a receive-only reference to the message channel.
// The message channel will be closed once all operations on this group have completed.
//
// The Wait function may only be called once.
func (op *OpGroup[M]) Wait() <-chan M {
	go func() {
		op.wg.Wait()
		close(op.c)
	}()
	return op.c
}
