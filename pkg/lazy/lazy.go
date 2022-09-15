package lazy

import "sync"

// Lazy is an object that a lazily-initialized value of type T.
//
// A Lazy must not be copied after first use.
type Lazy[T any] struct {
	once  sync.Once
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
	lazy.once.Do(func() {
		lazy.value = init()
	})
	return lazy.value
}
