// Package liquid provides Liquid
//
//spellchecker:words liquid
package liquid

//spellchecker:words github wisski distillery internal component instances malt models
import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/malt"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sqldelegator"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Liquid is the core of a WissKI Instance and used in every ingredient.
type Liquid struct {
	*malt.Malt
	models.Instance // TODO: move this into an explicit field

	DrupalUsername string
	DrupalPassword string
}

func (liquid *Liquid) DelegatedSQL() sqldelegator.DelegatedSQL {
	// Note: We cannot cache here, as the implementation itself might switch between global and local sql
	// depending on the configuration.
	return liquid.Delegator.For(liquid.Instance)
}
