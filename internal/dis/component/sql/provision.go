package sql

//spellchecker:words context errors strings time github wisski distillery internal models pkglib errorsx stream timex
import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/stream"
	"go.tkw01536.de/pkglib/timex"
)

// Provision provisions sql-specific resource for the given instance.
func (sql *SQL) Provision(ctx context.Context, instance models.Instance, domain string) error {
	return sql.CreateDatabase(ctx, CreateOpts{
		Name: instance.SqlDatabase,

		CreateUser: true,
		Username:   instance.SqlUsername,
		Password:   instance.SqlPassword,
	})
}

// Purge purges sql-specific resources for the given instance.
func (sql *SQL) Purge(ctx context.Context, instance models.Instance, domain string) error {
	return errorsx.Combine(
		sql.DropDatabase(ctx, instance.SqlDatabase),
		sql.DropUser(ctx, instance.SqlUsername),
	)
}

var errCreateSuperuserGrant = errors.New("`CreateSuperUser': grant failed")

// CreateSuperuser createsa new user, with the name 'user' and the password 'password'.
// CreateSuperuser always waits for the database to become available, and then uses the internal 'mysql' executable of the container.
func (sql *SQL) CreateSuperuser(ctx context.Context, user, password string, allowExisting bool) (e error) {
	stack, err := sql.OpenStack()
	if err != nil {
		return err
	}
	defer errorsx.Close(stack, &e, "stack")

	nilStream := stream.FromNil()

	// wait to connect to the databse and for the 'select 1' query to succeed
	if err := timex.TickUntilFunc(func(time.Time) bool {
		running, err := stack.Running(ctx)
		if err != nil || !running {
			return false
		}

		code := sql.Shell(ctx, nilStream, "-e", "select 1;")
		return code == 0
	}, ctx, sql.PollInterval); err != nil {
		return fmt.Errorf("failed to wait for sql: %w", err)
	}

	var IfNotExists string
	if allowExisting {
		IfNotExists = "IF NOT EXISTS"
	}

	var (
		userQuoted = quoteSingle(user)
		passQuoted = quoteSingle(password)
	)

	var builder strings.Builder
	code := sql.Shell(
		ctx, stream.NewIOStream(nil, &builder, nil), "-e",
		"CREATE USER "+IfNotExists+" "+userQuoted+"@'%' IDENTIFIED BY "+passQuoted+";"+
			"GRANT ALL PRIVILEGES ON *.* TO "+userQuoted+"@'%' WITH GRANT OPTION;"+
			"FLUSH PRIVILEGES;",
	)
	if code != 0 {
		return fmt.Errorf("%w: %s", errCreateSuperuserGrant, builder.String())
	}
	return nil
}
