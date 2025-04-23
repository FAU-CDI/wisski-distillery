//spellchecker:words barrel
package barrel

//spellchecker:words context time github wisski distillery internal component meta status ingredient locker mstore
import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
)

// Build builds or rebuilds the barrel connected to this instance.
//
// It also logs the current time into the metadata belonging to this instance.
func (barrel *Barrel) Build(ctx context.Context, progress io.Writer, start bool) error {
	if !barrel.dependencies.Locker.TryLock(ctx) {
		return locker.ErrLocked
	}
	defer barrel.dependencies.Locker.Unlock(ctx)

	stack := barrel.Stack()

	var context component.InstallationContext

	{
		err := stack.Install(ctx, progress, context)
		if err != nil {
			return err
		}
	}

	{
		err := stack.Update(ctx, progress, start)
		if err != nil {
			return err
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
		return t, err
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (barrel *Barrel) setLastRebuild(ctx context.Context) error {
	return lastRebuild.Set(ctx, barrel.dependencies.MStore, time.Now().Unix())
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
