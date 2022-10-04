package rlock

import (
	"sync"
	"time"
)

// RLock is like [sync.Mutex], but permits recursive locking.
type RLock struct {
	m sync.Mutex // m is held internally

	held    bool
	holder  int
	counter uint64
}

// Lock acquires this lock with the given id, and blocks until it can be aquired.
// Concurrent locks with the same ids do not block; however each should be unlocked with a call to unlock.
func (rm *RLock) Lock(id int) {
loop:
	for {
		rm.m.Lock()
		switch {
		case !rm.held:
			rm.held = true
			rm.holder = id
			break loop
		case rm.held && rm.holder == id:
			break loop
		}
		rm.m.Unlock()
		time.Sleep(time.Millisecond) // spinning!
	}

	rm.counter++
	rm.m.Unlock()
}

// Unlock releases the lock
func (rm *RLock) Unlock() {
	rm.m.Lock()
	defer rm.m.Unlock()

	if !rm.held || rm.counter <= 0 {
		panic("RLock: Unlock() without Lock()")
	}

	rm.counter--
	if rm.counter == 0 {
		rm.held = false
		rm.holder = 0
	}
}
