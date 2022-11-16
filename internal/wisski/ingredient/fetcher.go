package ingredient

import (
	"fmt"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
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
	Server *phpx.Server
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

	// Statistics of the wisski (TODO: fix me)
	Statistics Statistics

	// List of backups made
	Snapshots []models.Export

	// List of SSH Keys
	SSHKeys []string

	// WissKI content information
	NoPrefixes   bool              // TODO: Move this into the database
	Prefixes     []string          // list of prefixes
	Pathbuilders map[string]string // all the pathbuilders
}

type Statistics struct {
	Activity struct {
		MostVisited string `json:"mostVisited"`
		PageVisits  []struct {
			URL    string `json:"url"`
			Visits int    `json:"visits"`
		} `json:"pageVisits"`
		TotalEditsLastWeek int `json:"totalEditsLastWeek"`
	} `json:"activity"`
	Bundles     BundleStatistics `json:"bundles"`
	Triplestore struct {
		Graphs []struct {
			URI   string `json:"uri"`
			Count int    `json:"triples"`
		} `json:"graphStatistics"`
		Total int `json:"totalTriples"`
	} `json:"triplestore"`
	Users struct {
		LastLogin  string `json:"lastLogin"`
		TotalUsers int    `json:"totalUsers"`
	} `json:"users"`
}

type BundleStatistics struct {
	Bundles []struct {
		Label       string `json:"label"`
		MachineName string `json:"machineName"`

		Count int `json:"entities"`

		LastEdit int `json:"lastEdit"`

		MainBundle phpx.BooleanIsh `json:"mainBundle"`
	} `json:"bundleStatistics"`
	TotalBundles     int `json:"totalBundles"`
	TotalMainBundles int `json:"totalMainBundles"`
}

func (bs BundleStatistics) Summary() string {
	var totalCount int
	for _, bundle := range bs.Bundles {
		totalCount += bundle.Count
	}
	if totalCount == 0 {
		return ""
	}

	entitySubject := "Entities"
	if totalCount == 1 {
		entitySubject = "Entity"
	}

	bundleSubject := "Bundles"
	if len(bs.Bundles) == 1 {
		bundleSubject = "Bundle"
	}

	return fmt.Sprintf("%d %s in %d %s", totalCount, entitySubject, len(bs.Bundles), bundleSubject)
}
