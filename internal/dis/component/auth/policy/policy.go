//spellchecker:words policy
package policy

//spellchecker:words context github wisski distillery internal component auth models gorm
import (
	"context"
	"fmt"

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
		Model: models.Grant{},
	}
}

func (pol *Policy) openInterface(ctx context.Context) (gorm.Interface[models.Grant], error) {
	table, err := sql.OpenInterface[models.Grant](ctx, pol.dependencies.SQL, pol)
	if err != nil {
		return nil, fmt.Errorf("failed to open interface: %w", err)
	}
	return table, nil
}

func (pol *Policy) openDB(ctx context.Context) (*gorm.DB, error) {
	conn, err := pol.dependencies.SQL.OpenTable(ctx, pol)
	if err != nil {
		return nil, fmt.Errorf("failed to open interface: %w", err)
	}
	return conn, nil
}
