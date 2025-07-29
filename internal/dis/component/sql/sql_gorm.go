package sql

//spellchecker:words context errors reflect gorm driver mysql logger github wisski distillery internal component
import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

// OpenTable opens a *gorm.DB connection to the given table.
// Should use [OpenInterface] where possible.
//
// TODO: Migrate everything here!
func (sql *SQL) OpenTable(ctx context.Context, table component.Table) (*gorm.DB, error) {
	db, err := sql.connectGorm(ctx)
	if err != nil {
		return nil, err
	}

	name := table.TableInfo().Name()

	db = db.Table(name)
	if db.Error != nil {
		return nil, fmt.Errorf("failed to open connection to table %q: %w", name, db.Error)
	}

	return db, nil
}

var errWrongGenericType = errors.New("wrong generic type for table")

// OpenInterface opens a [gorm.Interface] to the given sql and table interface.
// The generic parameter T must correspond to the [component.Table]'s TableInfo.
func OpenInterface[T any](ctx context.Context, sql *SQL, table component.Table) (gorm.Interface[T], error) {
	info := table.TableInfo()
	if got := reflect.TypeFor[T](); got != reflect.TypeOf(info.Model) {
		return nil, fmt.Errorf("%w: got %v, expected %v", errWrongGenericType, got, info.Model)
	}

	db, err := sql.connectGorm(ctx)
	if err != nil {
		return nil, err
	}

	return gorm.G[T](db), nil
}

// creates a fresh gorm connection.
func (sql *SQL) connectGorm(ctx context.Context) (*gorm.DB, error) {
	conn, err := sql.connectSQL(ctx)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn: conn,

			DefaultStringSize: 256,
		}),
		&gorm.Config{
			Logger: newGormLogger().LogMode(logger.Info),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with gorm: %w", err)
	}

	db = db.WithContext(ctx)

	if db.Error != nil {
		return nil, db.Error
	}
	return db, nil
}
