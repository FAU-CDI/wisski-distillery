package extras

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"

	_ "embed"
)

// Prefixes implements reading and writing prefix
type Adapters struct {
	ingredient.Base
	Dependencies struct {
		PHP *php.PHP
	}
}

//go:embed adapters.php
var adaptersPHP string

type DistilleryAdapter struct {
	Label          string
	MachineName    string
	Description    string
	InstanceDomain string

	GraphDBRepository string
	GraphDBUsername   string
	GraphDBPassword   string
}

func (wisski *Adapters) CreateDistilleryAdapter(ctx context.Context, server *phpx.Server, adapter DistilleryAdapter) error {
	return wisski.Dependencies.PHP.ExecScript(
		ctx, server, nil, adaptersPHP,
		"create_distillery_adapter",
		adapter.Label, adapter.MachineName, adapter.Description, adapter.InstanceDomain, adapter.GraphDBRepository, adapter.GraphDBUsername, adapter.GraphDBPassword,
	)
}
