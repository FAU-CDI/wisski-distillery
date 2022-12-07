package cron

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/timex"
	"github.com/rs/zerolog"
)

type Cron struct {
	component.Base

	Tasks []component.Cronable
}

// Listen returns a channel that listens for triggers in the current process.
// It is intended to be passed to Start.
func (control *Cron) Listen(ctx context.Context) (<-chan struct{}, func()) {
	var (
		signals = make(chan os.Signal, 1)
		notify  = make(chan struct{}, 1)
	)

	go func() {
		for {
			select {
			case <-signals:
				notify <- struct{}{}
			case <-ctx.Done():
				return
			}
		}
	}()

	signal.Notify(signals, syscall.SIGHUP)
	return notify, func() {
		signal.Ignore(syscall.SIGHUP)
	}
}

// Once immediatly runs all cron jobs in the current thread.
// Once returns once all cron jobs have returned.
//
// Once should not be called concurrently with Cron.
func (control *Cron) Once(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(control.Tasks))

	zerolog.Ctx(ctx).Info().Time("time", time.Now()).Msg("Starting Cron")

	for _, task := range control.Tasks {
		go func(task component.Cronable) {
			defer wg.Done()

			name := task.TaskName()

			start := time.Now()
			zerolog.Ctx(ctx).Info().Str("task", name).Time("time", start).Msg("Calling Cron()")

			panicked, panik, err := func() (panicked bool, panik any, err error) {
				defer func() {
					panik = recover()
				}()

				panicked = true
				err = task.Cron(ctx)
				panicked = false

				return
			}()

			took := time.Since(start)

			switch {
			case !panicked:
				zerolog.Ctx(ctx).Err(err).Str("task", name).Dur("took", took).Msg("Finished Cron()")
			case panicked:
				zerolog.Ctx(ctx).Error().Str("task", name).Dur("took", took).Str("panic", fmt.Sprint(panik)).Msg("Finished Cron()")
			}
		}(task)
	}

	wg.Wait()
	zerolog.Ctx(ctx).Info().Time("time", time.Now()).Msg("Finished Cron")
}

// Start invokes all cron jobs regularly, waiting interval between invocations.
//
// The first run is invoked immediatly.
// The call to Start returns after the first invocation of all cron tasks.
//
// The returned channel is closed once no more cron tasks are active.
func (control *Cron) Start(ctx context.Context, interval time.Duration, signal <-chan struct{}) <-chan struct{} {
	// run runs cron tasks with the configured timeout
	run := func() {
		ctx, done := context.WithTimeout(ctx, interval)
		defer done()

		control.Once(ctx)
	}

	cleanup := make(chan struct{}) // closed once we have finished running everything

	run() // run tasks immediatly

	// start a new xgoroutine to run cron tasks
	go func() {
		defer close(cleanup)

		timer := timex.NewTimer()
		for {
			timex.StopTimer(timer)
			timer.Reset(interval)

			select {
			case <-timer.C:
				zerolog.Ctx(ctx).Debug().Msg("Cron() timer fired")
			case <-signal:
				zerolog.Ctx(ctx).Debug().Msg("Cron() received signal")
			case <-ctx.Done():
				timex.StopTimer(timer)
				return
			}

			run()
		}
	}()

	// and return the cleanup channel
	return cleanup
}
