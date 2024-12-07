//spellchecker:words purger
package purger

//spellchecker:words context github wisski distillery internal component instances models logging goprogram exit
import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Purger purges instances from the distillery
type Purger struct {
	component.Base
	dependencies struct {
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
	logging.LogMessage(out, "Checking bookkeeping table")
	instance, err := purger.dependencies.Instances.WissKI(ctx, slug)
	if err == instances.ErrWissKINotFound {
		fmt.Fprintln(out, "Not found in bookkeeping table, assuming defaults")
		instance, err = purger.dependencies.Instances.Create(slug, models.System{})
	}
	if err != nil {
		return errPurgeNoDetails.WithMessageF(err)
	}

	// remove docker stack
	logging.LogMessage(out, "Stopping and removing docker container")
	if err := instance.Barrel().Stack().Down(ctx, out); err != nil {
		fmt.Fprintln(out, err)
	}

	// remove the filesystem
	logging.LogMessage(out, "Removing from filesystem %s", instance.FilesystemBase)
	if err := os.RemoveAll(instance.FilesystemBase); err != nil {
		fmt.Fprintln(out, err)
	}

	// purge all the instance specific resources
	if err := logging.LogOperation(func() error {
		domain := instance.Domain()
		for _, pc := range purger.dependencies.Provisionable {
			logging.LogMessage(out, "Purging %s resources", pc.Name())
			err := pc.Purge(ctx, instance.Instance, domain)
			if err != nil {
				return err
			}
		}

		return nil
	}, out, "Purging instance-specific resources"); err != nil {
		return errPurgeGeneric.WithMessageF(slug, err)
	}

	// remove from bookkeeping
	logging.LogMessage(out, "Removing instance from bookkeeping")
	if err := instance.Bookkeeping().Delete(ctx); err != nil {
		fmt.Fprintln(out, err)
	}

	// remove the filesystem
	logging.LogMessage(out, "Remove lock data")
	if instance.Locker().TryUnlock(ctx) {
		fmt.Fprintln(out, "instance was not locked")
	}

	return nil
}
