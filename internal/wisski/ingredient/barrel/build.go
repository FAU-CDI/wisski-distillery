//spellchecker:words barrel
package barrel

//spellchecker:words context errors time github wisski distillery internal component meta status ingredient mstore pkglib errorsx
import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"go.tkw01536.de/pkglib/errorsx"
)

// Build builds or rebuilds the barrel connected to this instance.
//
// It also logs the current time into the metadata belonging to this instance.
func (barrel *Barrel) Build(ctx context.Context, progress io.Writer, start bool) (e error) {
	if err := barrel.dependencies.Locker.TryLock(ctx); err != nil {
		return fmt.Errorf("unable to lock instance: %w", err)
	}
	defer barrel.dependencies.Locker.Unlock(ctx)

	stack, err := barrel.OpenStack()
	if err != nil {
		return fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	var context component.InstallationContext

	{
		err := stack.Install(ctx, progress, context)
		if err != nil {
			return fmt.Errorf("failed to install stack: %w", err)
		}
	}

	{
		err := stack.Update(ctx, progress, start)
		if err != nil {
			return fmt.Errorf("failed to update stack: %w", err)
		}
	}

	// store the current last rebuild
	return barrel.setLastRebuild(ctx)
}

// TODO: Move this to time.Time.
var lastRebuild = mstore.For[int64]("lastRebuild")

func (barrel Barrel) LastRebuild(ctx context.Context) (t time.Time, err error) {
	epoch, err := lastRebuild.Get(ctx, barrel.dependencies.MStore)
	if errors.Is(err, meta.ErrMetadatumNotSet) {
		return t, nil
	}
	if err != nil {
		return t, fmt.Errorf("failed to set last rebuild: %w", err)
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (barrel *Barrel) setLastRebuild(ctx context.Context) error {
	if err := lastRebuild.Set(ctx, barrel.dependencies.MStore, time.Now().Unix()); err != nil {
		return fmt.Errorf("failed to set last rebuild: %w", err)
	}
	return nil
}

type LastRebuildFetcher struct {
	ingredient.Base
	dependencies struct {
		Barrel *Barrel
	}
}

var (
	_ ingredient.WissKIFetcher = (*LastRebuildFetcher)(nil)
)

func (lbr *LastRebuildFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.LastRebuild, _ = lbr.dependencies.Barrel.LastRebuild(flags.Context)
	return
}
