package ingredient

import (
	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
)

type WissKIFetcher interface {
	Ingredient

	// Fetch fetches information from this ingredient and writes it into target.
	// Distinct WissKIFetchers must write into distinct fields.
	Fetch(flags FetcherFlags, target *status.WissKI) error
}

// FetcherFlags describes options for a WissKIFetcher
type FetcherFlags struct {
	Quick  bool
	Server *phpx.Server
}
