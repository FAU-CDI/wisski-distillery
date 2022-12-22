package barrel

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// Running checks if this WissKI is currently running.
func (barrel *Barrel) Running(ctx context.Context, progress io.Writer) (bool, error) {
	ps, err := barrel.Stack().Ps(ctx, progress)
	if err != nil {
		return false, err
	}
	return len(ps) > 0, nil
}

type RunningFetcher struct {
	ingredient.Base
	Dependencies struct {
		Barrel *Barrel
	}
}

var (
	_ ingredient.WissKIFetcher = (*RunningFetcher)(nil)
)

func (rf *RunningFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.Running, err = rf.Dependencies.Barrel.Running(flags.Context, io.Discard)
	return
}
