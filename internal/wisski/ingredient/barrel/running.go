package barrel

import (
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/tkw1536/goprogram/stream"
)

// Running checks if this WissKI is currently running.
func (barrel *Barrel) Running() (bool, error) {
	ps, err := barrel.Stack().Ps(stream.FromNil())
	if err != nil {
		return false, err
	}
	return len(ps) > 0, nil
}

type RunningFetcher struct {
	ingredient.Base

	Barrel *Barrel
}

func (rf *RunningFetcher) Fetch(flags ingredient.FetchFlags, info *ingredient.Information) (err error) {
	info.Running, err = rf.Barrel.Running()
	return
}