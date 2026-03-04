package sqldelegator

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"go.tkw01536.de/pkglib/stream"
)

// Delegator delegates SQL operations for a specific instance
type Delegator struct {
	component.Base
	dependencies struct {
		SQL *sql.SQL
	}
}

type DelegatedSQL interface {
	// SQLUrl returns the SQL URL for the given instance.
	SQLUrl() string

	// Provision or purge databases belonging to a given instance.
	Provision(ctx context.Context) error
	Purge(ctx context.Context) error

	// Shell opens a shell inside the given sql database.
	Shell(ctx context.Context, io stream.IOStream, argv ...string) int

	// Snapshot makes a snapshot of the entire database.
	Snapshot(ctx context.Context, progress io.Writer, dest io.Writer) error
	// Restore restore the database from the given reader.
	Restore(ctx context.Context, reader io.Reader, io stream.IOStream) error
}

type delegated struct {
	*Impl
	instance  models.Instance
	delegator *Delegator
}

// For returns a new InstanceSQL for the given instance.
func (delegator *Delegator) For(instance models.Instance) DelegatedSQL {
	return &delegated{
		instance:  instance,
		delegator: delegator,

		Impl: NewImpl("sql", func() (*dockerx.Stack, error) {
			stack, err := delegator.dependencies.SQL.OpenStack()
			if err != nil {
				return nil, err
			}
			return stack.Stack, nil
		}),
	}
}
