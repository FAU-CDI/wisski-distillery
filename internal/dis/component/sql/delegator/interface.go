package delegator

import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// InstanceInterface implements all global component hooks for the delegator.
type InstanceInterface struct {
	component.Base
	dependencies struct {
		SQL *sql.SQL
	}
}

var (
	_ component.Provisionable = (*InstanceInterface)(nil)
	_ component.Snapshotable  = (*InstanceInterface)(nil)
)

func (iface *InstanceInterface) ProvisionNeedsStack(instance models.Instance) bool {
	return true
}

func (iface *InstanceInterface) Provision(ctx context.Context, instance models.Instance, domain string, stack *component.StackWithResources) error {
	return iface.dependencies.SQL.For(instance).Provision(ctx)
}

func (iface *InstanceInterface) Purge(ctx context.Context, instance models.Instance, domain string) error {
	return iface.dependencies.SQL.For(instance).Purge(ctx)
}

func (*InstanceInterface) SnapshotNeedsRunning(wisski models.Instance) bool { return false }

func (*InstanceInterface) SnapshotName() string { return "sql" }

func (iface *InstanceInterface) Snapshot(instance models.Instance, scontext *component.StagingContext) error {
	delegated := iface.dependencies.SQL.For(instance)
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
