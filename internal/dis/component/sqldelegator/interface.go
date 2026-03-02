package sqldelegator

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Interface implements all global component hooks for the delegator.
type Interface struct {
	component.Base
	dependencies struct {
		delegator *Delegator
	}
}

var (
	_ component.Provisionable = (*Interface)(nil)
)

func (iface *Interface) Provision(ctx context.Context, instance models.Instance, domain string) error {
	return iface.dependencies.delegator.For(instance).Provision(ctx)
}

func (iface *Interface) Purge(ctx context.Context, instance models.Instance, domain string) error {
	return iface.dependencies.delegator.For(instance).Purge(ctx)
}
