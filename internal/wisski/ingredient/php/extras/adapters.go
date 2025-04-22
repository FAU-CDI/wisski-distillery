//spellchecker:words extras
package extras

//spellchecker:words context github wisski distillery internal phpx ingredient embed
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"

	_ "embed"
)

// Prefixes implements reading and writing prefix.
type Adapters struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

//go:embed adapters.php
var adaptersPHP string

type DistilleryAdapter struct {
	ID string

	Label          string
	Description    string
	InstanceDomain string

	GraphDBRepository string
	GraphDBUsername   string
	GraphDBPassword   string
}

// Adapters returns a list of (managed) adapters belonging to the given WissKI.
func (wisski *Adapters) Adapters() []DistilleryAdapter {
	return []DistilleryAdapter{
		wisski.DefaultAdapter(),
	}
}

func (wisski *Adapters) DefaultAdapter() DistilleryAdapter {
	liquid := ingredient.GetLiquid(wisski)
	return DistilleryAdapter{
		ID: "default",

		Label:             "Default WissKI Distillery Adapter",
		Description:       "Default Adapter for " + liquid.Domain(),
		InstanceDomain:    liquid.Domain(),
		GraphDBRepository: liquid.GraphDBRepository,
		GraphDBUsername:   liquid.GraphDBUsername,
		GraphDBPassword:   liquid.GraphDBPassword,
	}
}

// SetAdapter creates or updates an adapter in the distillery.
// created indicates if a new adapter was created or if an existing one was updated.
func (wisski *Adapters) SetAdapter(ctx context.Context, server *phpx.Server, adapter DistilleryAdapter) (created bool, err error) {
	err = wisski.dependencies.PHP.ExecScript(
		ctx, server, &created, adaptersPHP,
		"create_or_update_distillery_adapter",
		adapter.Label, adapter.ID, adapter.Description, adapter.InstanceDomain, adapter.GraphDBRepository, adapter.GraphDBUsername, adapter.GraphDBPassword,
	)
	return
}
