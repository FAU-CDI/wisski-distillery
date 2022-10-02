package rlock

import (
	"sync"
	"time"
)

type RLock struct {
	m sync.Mutex // m is held internally

	held    bool
	holder  int
	counter uint64
}

func (rm *RLock) Lock(id int) {
	for {
		rm.m.Lock()
		if !rm.held {
			rm.held = true
			rm.holder = id
			break
		} else if rm.held && rm.holder == id {
			break
		} else {
			rm.m.Unlock()
			time.Sleep(time.Millisecond)
			continue
		}
	}

	rm.counter++
	rm.m.Unlock()
}

func (rm *RLock) Unlock() {
	rm.m.Lock()
	rm.counter--
	if rm.counter == 0 {
		rm.held = false
		rm.holder = 0
	}
	rm.m.Unlock()
}
