//spellchecker:words barrel
package barrel

//spellchecker:words context github wisski distillery internal status ingredient pkglib errorsx
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"go.tkw01536.de/pkglib/errorsx"
)

// Running checks if this WissKI is currently running.
func (barrel *Barrel) Running(ctx context.Context) (r bool, e error) {
	stack, err := barrel.OpenStack()
	if err != nil {
		return false, fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	containers, err := stack.Containers(ctx, false)
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
