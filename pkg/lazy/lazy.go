package lazy

import (
	"sync"
)

// Lazy holds a lazily initialized value of T.
// A non-zero lazy must not be copied after first use.
type Lazy[T any] struct {
	once sync.Once

	m     sync.RWMutex // m protects setting the value of this T
	value T            // the stored value
}

// Get returns the value associated with this Lazy.
//
// If no other call to Get has started or completed an initialization, calls init to initialize the value.
// A nil init function indicates to store the zero value of T.
// If an initialization has been previously completed, the previously stored value is returned.
//
// If init panics, the initization is considered to be completed.
// Future calls to Get() do not invoke init, and the zero value of T is returned.
//
// Get may safely be called concurrently.
func (lazy *Lazy[T]) Get(init func() T) T {
	lazy.m.RLock()
	defer lazy.m.RUnlock()

	lazy.once.Do(func() {
		if init != nil {
			lazy.value = init()
		}
	})

	return lazy.value
}

// Set atomically sets the value of this lazy.
// Any previously set value will be overwritten.
// Future calls to [Get] will not invoke init.
//
// It may be called concurrently with calls to [Get].
func (lazy *Lazy[T]) Set(value T) {
	lazy.m.Lock()
	defer lazy.m.Unlock()

	lazy.value = value
	lazy.once.Do(func() {})
}
