package policy

import (
	"context"
	"reflect"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"gorm.io/gorm"
)

type Policy struct {
	component.Base

	dependencies struct {
		SQL  *sql.SQL
		Auth *auth.Auth
	}
}

var (
	_ component.Provisionable  = (*Policy)(nil)
	_ component.UserDeleteHook = (*Policy)(nil)
	_ component.Table          = (*Policy)(nil)
)

func (pol *Policy) TableInfo() component.TableInfo {
	return component.TableInfo{
		Name:  models.GrantTable,
		Model: reflect.TypeFor[models.Grant](),
	}
}

func (pol *Policy) table(ctx context.Context) (*gorm.DB, error) {
	return pol.dependencies.SQL.QueryTable(ctx, pol)
}
