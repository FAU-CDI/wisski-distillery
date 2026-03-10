package impl

import (
	"context"
	"errors"
	"fmt"
	"io"

	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/stream"
)

// Bound represents an SQL implementation bound to a specific database, username and password.
type Bound struct {
	Impl Impl

	Username string
	Password string
	Database string
}

// Returns the database URL to be used for this bound sql instance.
func (bound *Bound) SQLUrl() string {
	return bound.Impl.URL(bound.Username, bound.Password, bound.Database)
}

// Provision provisions a new database for the given instance.
// It ensures that the database container is started and responding to queries afterwards.
func (bound *Bound) Provision(ctx context.Context) error {
	if err := bound.Impl.StartAndWait(ctx, stream.Null); err != nil {
		return fmt.Errorf("failed to start and wait for database: %w", err)
	}

	return bound.Impl.CreateDatabase(ctx, stream.Null, CreateOpts{
		Name: bound.Database,

		CreateUser: true,
		Username:   bound.Username,
		Password:   bound.Password,
	})
}

// Purge purges the database for the given instance.
func (bound *Bound) Purge(ctx context.Context) error {
	return bound.Impl.Purge(ctx, stream.Null, bound.Database, bound.Username)
}

// Shell opens a shell inside the given sql database.
func (bound *Bound) Shell(ctx context.Context, io stream.IOStream, argv ...string) int {
	return bound.Impl.Shell(ctx, io, argv...)
}

// Snapshot makes a snapshot of the database.
func (bound *Bound) Snapshot(ctx context.Context, progress io.Writer, dest io.Writer) error {
	return bound.Impl.SnapshotDB(ctx, progress, dest, bound.Database)
}

var errFailedToRestoreSQLContents = errors.New("failed to restore SQL contents")

// Restore restore the given database from the given reader.
// The database name is replaced inside the SQL dump.
func (bound *Bound) Restore(ctx context.Context, reader io.Reader, io stream.IOStream) (e error) {
	replacedFile := replaceSqlDatabaseName(reader, func(string) string { return bound.Database })
	defer errorsx.Close(replacedFile, &e, "replaced file")

	code := bound.Impl.Shell(ctx, stream.NewIOStream(io.Stdout, io.Stderr, replacedFile))
	if code != 0 {
		return fmt.Errorf("%w: exit code %d", errFailedToRestoreSQLContents, code)
	}
	return nil
}
