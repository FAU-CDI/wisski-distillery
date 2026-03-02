package sqldelegator

import (
	"context"
	"fmt"
	"io"

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
	_ component.Snapshotable  = (*Interface)(nil)
)

func (iface *Interface) Provision(ctx context.Context, instance models.Instance, domain string) error {
	return iface.dependencies.delegator.For(instance).Provision(ctx)
}

func (iface *Interface) Purge(ctx context.Context, instance models.Instance, domain string) error {
	return iface.dependencies.delegator.For(instance).Purge(ctx)
}

func (*Interface) SnapshotNeedsRunning() bool { return false }

func (*Interface) SnapshotName() string { return "sql" }

func (iface *Interface) Snapshot(instance models.Instance, scontext *component.StagingContext) error {
	delegated := iface.dependencies.delegator.For(instance)
	if err := scontext.AddDirectory(".", func(ctx context.Context) error {
		if err := scontext.AddFile("database.sql", func(ctx context.Context, file io.Writer) error {
			if err := delegated.Snapshot(ctx, scontext.Progress(), file); err != nil {
				return fmt.Errorf("failed to snapshot database: %w", err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to add sql file: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to snapshot directory: %w", err)
	}
	return nil
}
