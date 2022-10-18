// Package liquid provides Liquid
package liquid

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/malt"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Liquid is the core of a WissKI Instance and used in every ingredient.
type Liquid struct {
	*malt.Malt
	models.Instance

	DrupalUsername string
	DrupalPassword string
}
