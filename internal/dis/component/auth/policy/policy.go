package policy

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"gorm.io/gorm"
)

type Policy struct {
	component.Base

	Dependencies struct {
		SQL *sql.SQL
	}
}

var (
	_ component.Provisionable  = (*Policy)(nil)
	_ component.UserDeleteHook = (*Policy)(nil)
)

func (pol *Policy) table(ctx context.Context) (*gorm.DB, error) {
	return pol.Dependencies.SQL.QueryTable(ctx, true, models.GrantTable)
}
