package barrel

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/tkw1536/goprogram/stream"
)

// Running checks if this WissKI is currently running.
func (barrel *Barrel) Running(ctx context.Context) (bool, error) {
	ps, err := barrel.Stack().Ps(ctx, stream.FromNil())
	if err != nil {
		return false, err
	}
	return len(ps) > 0, nil
}

type RunningFetcher struct {
	ingredient.Base

	Barrel *Barrel
}

var (
	_ ingredient.WissKIFetcher = (*RunningFetcher)(nil)
)

func (rf *RunningFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.Running, err = rf.Barrel.Running(flags.Context)
	return
}
