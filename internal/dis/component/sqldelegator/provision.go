package sqldelegator

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"go.tkw01536.de/pkglib/errorsx"
)

func (delegator *delegated) Provision(ctx context.Context) error {
	return delegator.delegator.dependencies.SQL.CreateDatabase(ctx, sql.CreateOpts{
		Name: delegator.instance.SqlDatabase,

		CreateUser: true,
		Username:   delegator.instance.SqlUsername,
		Password:   delegator.instance.SqlPassword,
	})
}

func (delegator *delegated) Purge(ctx context.Context) error {
	return errorsx.Combine(
		delegator.delegator.dependencies.SQL.DropDatabase(ctx, delegator.instance.SqlDatabase),
		delegator.delegator.dependencies.SQL.DropUser(ctx, delegator.instance.SqlUsername),
	)
}
