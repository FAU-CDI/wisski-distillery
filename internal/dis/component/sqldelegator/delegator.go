package sqldelegator

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Delegator delegates SQL operations for a specific instance
type Delegator struct {
	component.Base
	dependencies struct {
		SQL *sql.SQL
	}
}

type DelegatedSQL interface {
	// Provision or purge databases belonging to a given instance.
	Provision(ctx context.Context) error
	Purge(ctx context.Context) error

	/*
		// Opens an sql shell for the given instance's sql database.
		Shell(ctx context.Context, io stream.IOStream) error

		// Snapshot or restore databases belonging to the given instance.
		Snapshot(ctx context.Context, progress io.Writer, dest io.Writer) error
		Restore(ctx context.Context, progress io.Writer, src io.Reader) error
	*/
}

type delegated struct {
	instance  models.Instance
	delegator *Delegator
}

// For returns a new InstanceSQL for the given instance.
func (delegator *Delegator) For(instance models.Instance) DelegatedSQL {
	return &delegated{
		instance:  instance,
		delegator: delegator,
	}
}
