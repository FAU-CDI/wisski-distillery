//spellchecker:words purger
package purger

//spellchecker:words context errors github wisski distillery internal component instances models logging pkglib errorsx
import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/errorsx"
)

// Purger purges instances from the distillery.
type Purger struct {
	component.Base
	dependencies struct {
		Instances     *instances.Instances
		Provisionable []component.Provisionable
	}
}

// Purge permanently purges an instance from the distillery.
// The instance does not have to exist; in which case the resources are also deleted.
func (purger *Purger) Purge(ctx context.Context, out io.Writer, slug string) (e error) {
	if _, err := logging.LogMessage(out, "Checking bookkeeping table"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	instance, err := purger.dependencies.Instances.WissKI(ctx, slug)
	if errors.Is(err, instances.ErrWissKINotFound) {
		_, _ = fmt.Fprintln(out, "Not found in bookkeeping table, assuming defaults")
		instance, err = purger.dependencies.Instances.Create(slug, models.System{})
	}
	if err != nil {
		return fmt.Errorf("unable to find instance details for purge: %w", err)
	}

	// remove docker stack
	if _, err := logging.LogMessage(out, "Stopping and removing docker container"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	stack, err := instance.Barrel().OpenStack()
	if err != nil {
		return fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	if err := stack.Down(ctx, out); err != nil {
		_, _ = fmt.Fprintln(out, err)
	}

	// remove the filesystem
	if _, err := logging.LogMessage(out, "Removing from filesystem %s", instance.FilesystemBase); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if err := os.RemoveAll(instance.FilesystemBase); err != nil {
		_, _ = fmt.Fprintln(out, err) // already handling error
	}

	// purge all the instance specific resources
	if err := logging.LogOperation(func() error {
		domain := instance.Domain()
		for _, pc := range purger.dependencies.Provisionable {
			if _, err := logging.LogMessage(out, "Purging %s resources", pc.Name()); err != nil {
				return fmt.Errorf("failed to log message: %w", err)
			}
			err := pc.Purge(ctx, instance.Instance, domain)
			if err != nil {
				return fmt.Errorf("failed to purge %s for instance %q: %w", pc.Name(), instance.Slug, err)
			}
		}

		return nil
	}, out, "Purging instance-specific resources"); err != nil {
		return fmt.Errorf("unable to purge instance %q: %w", slug, err)
	}

	// remove from bookkeeping
	if _, err := logging.LogMessage(out, "Removing instance from bookkeeping"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if err := instance.Bookkeeping().Delete(ctx); err != nil {
		_, _ = fmt.Fprintln(out, err)
	}

	// remove the filesystem
	if _, err := logging.LogMessage(out, "Remove lock data"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if instance.Locker().TryUnlock(ctx) {
		_, _ = fmt.Fprintln(out, "instance was not locked")
	}

	return nil
}
