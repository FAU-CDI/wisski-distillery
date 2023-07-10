package extras

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"golang.org/x/exp/slices"

	_ "embed"
)

type Requirements struct {
	ingredient.Base
	Dependencies struct {
		PHP *php.PHP
	}
}

var (
	_ ingredient.WissKIFetcher = (*Requirements)(nil)
)

//go:embed requirements.php
var requirementsPHP string

// Create creates a new block with the given title and html content
func (requirements *Requirements) Get(ctx context.Context, server *phpx.Server) (data []status.Requirement, err error) {
	err = requirements.Dependencies.PHP.ExecScript(ctx, server, &data, requirementsPHP, "get_requirements", requirements.URL().String())
	if err == nil {
		// sort first by weight, then by id!
		slices.SortFunc(data, func(a, b status.Requirement) bool {
			// compare first by weight
			if a.Weight < b.Weight {
				return true
			}
			if a.Weight > b.Weight {
				return false
			}

			// then by severity
			if a.Severity < b.Severity {
				return true
			}
			if a.Severity > b.Severity {
				return false
			}

			// and finally by id
			return a.ID < b.ID
		})
	}
	return
}

// Fetch fetches information
func (requirements *Requirements) Fetch(flags ingredient.FetcherFlags, target *status.WissKI) error {
	if flags.Quick {
		return nil
	}

	target.Requirements, _ = requirements.Get(flags.Context, flags.Server)
	return nil
}
