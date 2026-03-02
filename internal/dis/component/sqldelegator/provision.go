package sqldelegator

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"go.tkw01536.de/pkglib/errorsx"
)

func (delegated *delegated) SQLUrl() string {
	return "mysql://" + delegated.instance.SqlUsername + ":" + delegated.instance.SqlPassword + "@sql/" + delegated.instance.SqlDatabase
}

func (delegated *delegated) Provision(ctx context.Context) error {
	return delegated.delegator.dependencies.SQL.CreateDatabase(ctx, sql.CreateOpts{
		Name: delegated.instance.SqlDatabase,

		CreateUser: true,
		Username:   delegated.instance.SqlUsername,
		Password:   delegated.instance.SqlPassword,
	})
}

func (delegated *delegated) Purge(ctx context.Context) error {
	return errorsx.Combine(
		delegated.delegator.dependencies.SQL.DropDatabase(ctx, delegated.instance.SqlDatabase),
		delegated.delegator.dependencies.SQL.DropUser(ctx, delegated.instance.SqlUsername),
	)
}
