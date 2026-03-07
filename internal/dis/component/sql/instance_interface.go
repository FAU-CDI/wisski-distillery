package sql

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

var (
	_ component.Provisionable = (*SQL)(nil)
	_ component.Snapshotable  = (*SQL)(nil)
)

func (sql *SQL) ProvisionNeedsStack(instance models.Instance) bool {
	return instance.DedicatedSQL
}

var errFailedToProvision = errors.New("failed to provision sql database")

func (sql *SQL) Provision(ctx context.Context, instance models.Instance, domain string, stack *component.StackWithResources) error {
	provisionErr := sql.For(instance).Provision(ctx)
	if provisionErr == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", errFailedToProvision, provisionErr)
}

var errFailedToPurge = errors.New("failed to purge sql database")

func (sql *SQL) Purge(ctx context.Context, instance models.Instance, domain string) error {
	purgeErr := sql.For(instance).Purge(ctx)
	// ignore error while purging if we are using a dedicated sql server.
	// because it'll be deleted anyways by deleting the stack.
	if instance.DedicatedSQL {
		return nil
	}
	if purgeErr == nil {
		return nil
	}
	return fmt.Errorf("%w: %w", errFailedToPurge, purgeErr)
}

func (sql *SQL) SnapshotNeedsRunning(wisski models.Instance) bool { return false }

func (sql *SQL) SnapshotName() string { return "sql" }

func (sql *SQL) Snapshot(instance models.Instance, scontext *component.StagingContext) error {
	delegated := sql.For(instance)
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
