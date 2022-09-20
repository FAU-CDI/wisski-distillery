package component

import (
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Provisionable represents a component with a Provision method.
type Provisionable interface {
	Component

	// Provision provisions resources specific to the provided instance.
	// Domain holds the full (unique) domain name of the given instance.
	Provision(instance models.Instance, domain string) error
}
