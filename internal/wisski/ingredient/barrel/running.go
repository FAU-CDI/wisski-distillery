//spellchecker:words barrel
package barrel

//spellchecker:words context github wisski distillery internal status ingredient compose spec errdefs
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// Running checks if this WissKI is currently running.
func (barrel *Barrel) Running(ctx context.Context) (bool, error) {
	containers, err := barrel.Stack().Containers(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get containers: %w", err)
	}
	return len(containers) > 0, nil
}

type RunningFetcher struct {
	ingredient.Base
	dependencies struct {
		Barrel *Barrel
	}
}

var (
	_ ingredient.WissKIFetcher = (*RunningFetcher)(nil)
)

func (rf *RunningFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.Running, err = rf.dependencies.Barrel.Running(flags.Context)
	return
}
