package ingredient

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/phpserver"
)

// Fetcher is an ingredient with a fetch method
type Fetcher interface {
	Ingredient

	// Fetch fetchs information with the given information and writes it into info.
	// Distinct Fetchers must write into distinct fields.
	Fetch(flags FetchFlags, info *Information) error
}

// FetchFlags specifies what information to fetch
type FetchFlags struct {
	Quick  bool
	Server *phpserver.Server
}

// Information represents fetched information about a WissKI
type Information struct {
	Time time.Time // Time this info was built

	// Generic Information
	Slug string // slug
	URL  string // complete URL, including http(s)

	Locked bool // Is this instance currently locked?

	// Information about the running instance
	Running     bool
	LastRebuild time.Time
	LastUpdate  time.Time
	LastCron    time.Time

	// List of backups made
	Snapshots []models.Export

	// WissKI content information
	NoPrefixes   bool              // TODO: Move this into the database
	Prefixes     []string          // list of prefixes
	Pathbuilders map[string]string // all the pathbuilders
}
