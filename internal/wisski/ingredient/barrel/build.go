package barrel

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"github.com/tkw1536/goprogram/stream"
)

// Build builds or rebuilds the barel connected to this instance.
//
// It also logs the current time into the metadata belonging to this instance.
func (barrel *Barrel) Build(stream stream.IOStream, start bool) error {
	if !barrel.Locker.TryLock() {
		err := locker.Locked
		return err
	}
	defer barrel.Locker.Unlock()

	stack := barrel.Stack()

	var context component.InstallationContext

	{
		err := stack.Install(stream, context)
		if err != nil {
			return err
		}
	}

	{
		err := stack.Update(stream, start)
		if err != nil {
			return err
		}
	}

	// store the current last rebuild
	return barrel.setLastRebuild()
}

// TODO: Move this to time.Time
var lastRebuild = mstore.For[int64]("lastRebuild")

func (barrel Barrel) LastRebuild() (t time.Time, err error) {
	epoch, err := lastRebuild.Get(barrel.MStore)
	if err == meta.ErrMetadatumNotSet {
		return t, nil
	}
	if err != nil {
		return t, err
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (barrel *Barrel) setLastRebuild() error {
	return lastRebuild.Set(barrel.MStore, time.Now().Unix())
}
