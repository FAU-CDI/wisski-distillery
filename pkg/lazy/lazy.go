package lazy

import (
	"sync"
)

// Lazy is an object that a lazily-initialized value of type T.
//
// A Lazy must not be copied after first use.
type Lazy[T any] struct {
	once sync.Once

	m     sync.RWMutex // m protects setting the value of this T
	value T
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
