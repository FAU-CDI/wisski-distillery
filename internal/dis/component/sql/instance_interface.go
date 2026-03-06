package sql

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// InstanceSQL implements all global component hooks for the delegator.
type InstanceSQL struct {
	component.Base
	dependencies struct {
		SQL *SQL
	}
}

var (
	_ component.Provisionable = (*InstanceSQL)(nil)
	_ component.Snapshotable  = (*InstanceSQL)(nil)
)

func (iface *InstanceSQL) ProvisionNeedsStack(instance models.Instance) bool {
	return true
}

var errFailedToProvision = errors.New("failed to provision sql database")

func (iface *InstanceSQL) Provision(ctx context.Context, instance models.Instance, domain string, stack *component.StackWithResources) error {
	provisionErr := iface.dependencies.SQL.For(instance).Provision(ctx)
	return fmt.Errorf("%w: %w", errFailedToProvision, provisionErr)
}

var errFailedToPurge = errors.New("failed to purge sql database")

func (iface *InstanceSQL) Purge(ctx context.Context, instance models.Instance, domain string) error {
	purgeErr := iface.dependencies.SQL.For(instance).Purge(ctx)
	// ignore error while purging if we are using a dedicated sql server.
	// because it'll be deleted anyways by deleting the stack.
	if instance.DedicatedSQL {
		return nil
	}
	return fmt.Errorf("%w: %w", errFailedToPurge, purgeErr)
}

func (*InstanceSQL) SnapshotNeedsRunning(wisski models.Instance) bool { return false }

func (*InstanceSQL) SnapshotName() string { return "sql" }

func (iface *InstanceSQL) Snapshot(instance models.Instance, scontext *component.StagingContext) error {
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
