package component

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Observer is a component with an Observe method
type Observer interface {
	Component

	// Observe observes this distillery component and writes the result into observation
	// Distinct Observers must write into distinct fields.
	Observe(flags ObservationFlags, observation *Observation) error
}

type ObservationFlags struct{}

// Observation represents fetched information about the distillery
type Observation struct {
	Time time.Time // Time this obervation was built

	// Configuration of the distillery
	Config *config.Config

	// number of instances
	TotalCount   int
	RunningCount int
	StoppedCount int

	Backups []models.Export // list of backups
}
