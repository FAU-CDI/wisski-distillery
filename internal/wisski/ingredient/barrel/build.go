package barrel

import (
	"context"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
)

// Build builds or rebuilds the barel connected to this instance.
//
// It also logs the current time into the metadata belonging to this instance.
func (barrel *Barrel) Build(ctx context.Context, progress io.Writer, start bool) error {
	if !barrel.Dependencies.Locker.TryLock(ctx) {
		return locker.Locked
	}
	defer barrel.Dependencies.Locker.Unlock(ctx)

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

// TODO: Move this to time.Time
var lastRebuild = mstore.For[int64]("lastRebuild")

func (barrel Barrel) LastRebuild(ctx context.Context) (t time.Time, err error) {
	epoch, err := lastRebuild.Get(ctx, barrel.Dependencies.MStore)
	if err == meta.ErrMetadatumNotSet {
		return t, nil
	}
	if err != nil {
		return t, err
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (barrel *Barrel) setLastRebuild(ctx context.Context) error {
	return lastRebuild.Set(ctx, barrel.Dependencies.MStore, time.Now().Unix())
}

type LastRebuildFetcher struct {
	ingredient.Base
	Dependencies struct {
		Barrel *Barrel
	}
}

func (lbr *LastRebuildFetcher) Fetch(ctx context.Context, flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.LastRebuild, _ = lbr.Dependencies.Barrel.LastRebuild(ctx)
	return
}
