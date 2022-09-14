package wisski

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
)

// Distillery represents a WissKI Distillery
//
// It is the main structure used to interact with different components.
type Distillery struct {
	// Config holds the configuration of the distillery.
	// It is read directly from a configuration file.
	Config *config.Config

	// Upstream holds information to connect to the various running
	// distillery components.
	//
	// NOTE(twiesing): This is intended to eventually allow full remote management of the distillery.
	// But for now this will just hold upstream configuration.
	Upstream Upstream

	// components hold references to the various components of the distillery.
	components
}

// Upstream are the upstream urls connecting to the various external components.
type Upstream struct {
	SQL         string
	Triplestore string
}

// Context returns a new Context belonging to this distillery
func (dis *Distillery) Context() context.Context {
	return context.Background()
}
