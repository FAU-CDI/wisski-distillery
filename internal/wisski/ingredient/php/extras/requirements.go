//spellchecker:words extras
package extras

//spellchecker:words context strings github wisski distillery internal phpx status ingredient golang slices embed
import (
	"context"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"slices"

	_ "embed"
)

type Requirements struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

var (
	_ ingredient.WissKIFetcher = (*Requirements)(nil)
)

//go:embed requirements.php
var requirementsPHP string

// Create creates a new block with the given title and html content.
func (requirements *Requirements) Get(ctx context.Context, server *phpx.Server) (data []status.Requirement, err error) {
	err = requirements.dependencies.PHP.ExecScript(ctx, server, &data, requirementsPHP, "get_requirements", ingredient.GetLiquid(requirements).URL().String())
	if err == nil {
		// sort first by weight, then by id!
		slices.SortFunc(data, func(a, b status.Requirement) int {
			// compare first by severity
			if a.Severity < b.Severity {
				return 1
			}
			if a.Severity > b.Severity {
				return -1
			}

			// then by weight
			if a.Weight < b.Weight {
				return 1
			}
			if a.Weight > b.Weight {
				return -1
			}

			// and finally by id
			return strings.Compare(a.ID, b.ID)
		})
	}
	return
}

// Fetch fetches information.
func (requirements *Requirements) Fetch(flags ingredient.FetcherFlags, target *status.WissKI) error {
	if flags.Quick {
		return nil
	}

	target.Requirements, _ = requirements.Get(flags.Context, flags.Server)
	return nil
}
