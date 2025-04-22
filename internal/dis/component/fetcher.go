//spellchecker:words component
package component

//spellchecker:words context github wisski distillery internal status
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/status"
)

type DistilleryFetcher interface {
	Component

	// Fetch fetches information from this component and writes it into target.
	// Distinct DistilleryFetchers must write into distinct fields.
	Fetch(flags FetcherFlags, target *status.Distillery) error
}

// FetcherFlags describes options for a DistilleryFetcher.
type FetcherFlags struct {
	Context context.Context
}
