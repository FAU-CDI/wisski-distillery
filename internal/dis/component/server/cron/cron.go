//spellchecker:words cron
package cron

//spellchecker:words context signal sync syscall time github wisski distillery internal component wdlog pkglib timex
import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"go.tkw01536.de/pkglib/timex"
)

type Cron struct {
	component.Base
	dependencies struct {
		Tasks []component.Cronable
	}
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
	wg.Add(len(control.dependencies.Tasks))

	wdlog.Of(ctx).Info(
		"Starting Cron",
	)

	for _, task := range control.dependencies.Tasks {
		go func(task component.Cronable) {
			defer wg.Done()

			name := task.TaskName()

			start := time.Now()
			wdlog.Of(ctx).Info(
				"Calling Cron()",
				"task", name,
				"start", start,
			)

			panicked, panik, err := func() (panicked bool, panik any, err error) {
				defer func() {
					if panik = recover(); panik != nil {
						panicked = true
					}
				}()
				err = task.Cron(ctx)
				return
			}()

			took := time.Since(start)

			switch {
			case !panicked:
				if err != nil {
					wdlog.Of(ctx).Error(
						"Finished Cron()",
						"error", err,
						"task", name,
						"took", took,
					)
				} else {
					wdlog.Of(ctx).Info(
						"Finished Cron()",
						"task", name,
						"took", took,
					)
				}
			case panicked:
				wdlog.Of(ctx).Error(
					"Finished Cron()",
					"panic", fmt.Sprint(panik),
					"task", name,
					"took", took,
				)
			}
		}(task)
	}

	wg.Wait()
	wdlog.Of(ctx).Info(
		"Finished Cron",
	)
}

// Start invokes all cron jobs regularly, waiting between invocations as specified in configuration.
//
// A first run is invoked immediatly.
// The call to Start returns after the first invocation of all cron tasks.
//
// The returned channel is closed once no more cron tasks are active.
func (control *Cron) Start(ctx context.Context, signal <-chan struct{}) <-chan struct{} {
	interval := component.GetStill(control).Config.CronInterval
	wdlog.Of(ctx).Info(
		"Scheduling Cron() tasks",

		"interval", interval,
	)

	// run runs cron tasks with the configured timeout
	run := func() {
		ctx, done := context.WithTimeout(ctx, interval)
		defer done()

		control.Once(ctx)
	}

	cleanup := make(chan struct{}) // closed once we have finished running everything

	// start a new xgoroutine to run cron tasks
	go func() {
		defer close(cleanup)

		wdlog.Of(ctx).Debug("Cron() starting first run")
		run()
		wdlog.Of(ctx).Debug("Cron() beginnning scheduling")

		t := timex.NewTimer()
		defer timex.ReleaseTimer(t)
		for {
			timex.StopTimer(t)
			t.Reset(interval)

			select {
			case <-t.C:
				wdlog.Of(ctx).Debug("Cron() timer fired")
			case <-signal:
				wdlog.Of(ctx).Debug("Cron() received signal")
			case <-ctx.Done():
				return
			}

			run()
		}
	}()

	// and return the cleanup channel
	return cleanup
}
