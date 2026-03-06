package sql

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql/impl"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
)

// For returns a new BoundSQL for the given instance.
func (sql *SQL) For(instance models.Instance) *impl.Bound {

	var service string
	var openStack func() (*dockerx.Stack, error)
	if instance.DedicatedSQL {
		service = "dedicatedsql"
		openStack = func() (*dockerx.Stack, error) {
			return dockerx.NewStack(sql.dependencies.Docker, instance.FilesystemBase)
		}
	} else {
		service = "sql"
		openStack = func() (*dockerx.Stack, error) {
			stack, err := sql.OpenStack()
			if err != nil {
				return nil, err
			}
			return stack.Stack, nil
		}
	}

	return &impl.Bound{
		Username: instance.SqlUsername,
		Password: instance.SqlPassword,
		Database: instance.SqlDatabase,

		Impl: impl.New(service, openStack),
	}
}

func (sql *SQL) Global() *impl.Bound {
	config := component.GetStill(sql).Config.SQL

	return &impl.Bound{
		Username: config.AdminUsername,
		Password: config.AdminPassword,
		Database: config.Database,

		Impl: impl.New("sql", func() (*dockerx.Stack, error) {
			stack, err := sql.OpenStack()
			if err != nil {
				return nil, err
			}
			return stack.Stack, nil
		}),
	}
}
