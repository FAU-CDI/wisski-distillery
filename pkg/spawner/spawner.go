package spawner

import (
	"io"
	"log"
	"sync"
)

type Spawner[T io.Closer] struct {
	Spawn func() T
	Alive func(t T) bool

	Limit int

	initOnce sync.Once
	wg       sync.WaitGroup
	tasks    chan task[T]

	instances []T
	alive     []bool
}

type task[T any] struct {
	f    func(t T)
	done chan<- struct{}
}

func (tt task[T]) Do(t T) {
	defer close(tt.done)
	tt.f(t)
}

func (spawner *Spawner[T]) worker(i int) {
	spawner.wg.Add(1)
	go func() {
		defer spawner.wg.Done()

		var zero T
		for task := range spawner.tasks {
			// spawn a new instance once finished
			if !spawner.alive[i] {
				log.Println("respawning instance ", i)
				spawner.instances[i] = spawner.Spawn()
				spawner.alive[i] = true
			}

			task.Do(spawner.instances[i])

			// if the instance has died during execution
			// then do a close (to deallocate)
			if !spawner.Alive(spawner.instances[i]) {
				spawner.instances[i].Close()
				spawner.instances[i] = zero
				spawner.alive[i] = false
			}
		}
	}()
}

func (spawner *Spawner[T]) start() {
	spawner.initOnce.Do(func() {
		limit := spawner.Limit
		if limit < 1 {
			limit = 1
		}

		spawner.tasks = make(chan task[T], limit)
		spawner.alive = make([]bool, limit)
		spawner.instances = make([]T, limit)
		for i := 0; i < limit; i++ {
			spawner.worker(i)
		}
	})
}

// Do performs f on an unspecified object from the spawner.
// If no objects are available within time.Duration, a new object is spawned as there are at most limit objects.
func (spawner *Spawner[T]) Do(f func(t T)) {
	spawner.start()

	done := make(chan struct{})
	spawner.tasks <- task[T]{
		f:    f,
		done: done,
	}
	<-done
}

func (spawner *Spawner[T]) Close() {
	close(spawner.tasks)
	spawner.wg.Wait()

	spawner.wg.Add(len(spawner.alive))

	var zero T
	for i := range spawner.alive {
		go func(i int) {
			defer spawner.wg.Done()

			if spawner.alive[i] {
				spawner.instances[i].Close()
				spawner.instances[i] = zero
			}
		}(i)
	}

	spawner.wg.Wait()
}
