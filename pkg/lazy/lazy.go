package lazy

import (
	"sync"
	"time"
)

// Lazy is an object that a lazily-initialized value of type T.
//
// A Lazy must not be copied after first use.
type Lazy[T any] struct {
	once  sync.Once
	value T

	m         sync.RWMutex // m protects resetting this lazy
	lastReset time.Time    // last time this mutex was reset
}

// Get returns the value associated with this Lazy.
//
// If no other call to Get has started or completed an initialization, initializes the value using the init function.
// Otherwise, it returns the initialized value.
//
// If init panics, the initization is considered to be completed.
// Future calls to Get() do not invoke init, and the zero value of T is returned.
//
// Get may safely be called concurrently.
func (lazy *Lazy[T]) Get(init func() T) T {
	lazy.m.RLock()
	defer lazy.m.RUnlock()

	lazy.once.Do(func() {
		lazy.value = init()
	})
	return lazy.value
}

// Set atomically sets the value of this lazy, preventing future calls to get from invoking init.
// It may be called concurrently with calls to [Get] and [Reset].
func (lazy *Lazy[T]) Set(value T) {
	lazy.m.Lock()
	defer lazy.m.Unlock()

	lazy.value = value
	lazy.once.Do(func() {})
}

// Reset resets this Lazy, deleting any previously associated value.
//
// May be called concurrently with [Get].
// Future calls to [Get] will invoke init.
func (lazy *Lazy[T]) Reset() {
	lazy.m.Lock()
	defer lazy.m.Unlock()

	lazy.reset()
}

// ResetAfter resets this lazy if more than d time has passed since the last reset.
// If ResetAfter cannot lock, then it does not reset.
//
// May be called concurrently with other functions on this lazy.
func (lazy *Lazy[T]) ResetAfter(d time.Duration) {
	if !lazy.m.TryLock() {
		return
	}
	defer lazy.m.Unlock()

	if time.Since(lazy.lastReset) < d {
		return
	}

	lazy.reset()
}

// reset implements resetting this lazy.
// m myst be held for writing.
func (lazy *Lazy[T]) reset() {
	// reset the once
	lazy.once = sync.Once{}

	// reset the value
	var t T
	lazy.value = t

	// time of the last reset
	lazy.lastReset = time.Now()
}
