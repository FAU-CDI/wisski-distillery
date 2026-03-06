package delegator

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql/impl"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
)

// Delegator delegates SQL operations for a specific instance
type Delegator struct {
	component.Base
	dependencies struct {
		SQL *sql.SQL
	}
}

// Global returns a new BoundSQL for the global SQL database.
func (delegator *Delegator) Global() *impl.Bound {
	config := component.GetStill(delegator).Config.SQL

	return &impl.Bound{
		Username: config.AdminUsername,
		Password: config.AdminPassword,
		Database: config.Database,

		Impl: impl.New("sql", func() (*dockerx.Stack, error) {
			stack, err := delegator.dependencies.SQL.OpenStack()
			if err != nil {
				return nil, err
			}
			return stack.Stack, nil
		}),
	}
}
