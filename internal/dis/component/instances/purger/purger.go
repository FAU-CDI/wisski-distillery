package purger

import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Purger purges instances from the distillery
type Purger struct {
	component.Base
	Dependencies struct {
		Instances     *instances.Instances
		Provisionable []component.Provisionable
	}
}

var errPurgeNoDetails = exit.Error{
	Message:  "unable to find instance details for purge: %s",
	ExitCode: exit.ExitGeneric,
}
var errPurgeGeneric = exit.Error{
	Message:  "unable to purge instance %q: %s",
	ExitCode: exit.ExitGeneric,
}

// Purge permanently purges an instance from the distillery.
// The instance does not have to exist; in which case the resources are also deleted.
func (purger *Purger) Purge(ctx context.Context, out io.Writer, slug string) error {
	logging.LogMessage(out, ctx, "Checking bookkeeping table")
	instance, err := purger.Dependencies.Instances.WissKI(ctx, slug)
	if err == instances.ErrWissKINotFound {
		fmt.Fprintln(out, "Not found in bookkeeping table, assuming defaults")
		instance, err = purger.Dependencies.Instances.Create(slug)
	}
	if err != nil {
		return errPurgeNoDetails.WithMessageF(err)
	}

	// remove docker stack
	logging.LogMessage(out, ctx, "Stopping and removing docker container")
	if err := instance.Barrel().Stack().Down(ctx, out); err != nil {
		fmt.Fprintln(out, err)
	}

	// remove the filesystem
	logging.LogMessage(out, ctx, "Removing from filesystem %s", instance.FilesystemBase)
	if err := purger.Environment.RemoveAll(instance.FilesystemBase); err != nil {
		fmt.Fprintln(out, err)
	}

	// purge all the instance specific resources
	if err := logging.LogOperation(func() error {
		domain := instance.Domain()
		for _, pc := range purger.Dependencies.Provisionable {
			logging.LogMessage(out, ctx, "Purging %s resources", pc.Name())
			err := pc.Purge(ctx, instance.Instance, domain)
			if err != nil {
				return err
			}
		}

		return nil
	}, out, ctx, "Purging instance-specific resources"); err != nil {
		return errPurgeGeneric.WithMessageF(slug, err)
	}

	// remove from bookkeeping
	logging.LogMessage(out, ctx, "Removing instance from bookkeeping")
	if err := instance.Bookkeeping().Delete(ctx); err != nil {
		fmt.Fprintln(out, err)
	}

	// remove the filesystem
	logging.LogMessage(out, ctx, "Remove lock data")
	if instance.Locker().TryUnlock(ctx) {
		fmt.Fprintln(out, "instance was not locked")
	}

	return nil
}
