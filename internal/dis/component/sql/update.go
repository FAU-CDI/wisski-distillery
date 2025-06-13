package sql

//spellchecker:words context errors reflect time github wisski distillery internal component execx logging pkglib sqlx stream timex
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

// Update initializes or updates the SQL database.
func (sql *SQL) Update(ctx context.Context, progress io.Writer) error {
	{
		if _, err := logging.LogMessage(progress, "Creating administrative user"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		config := component.GetStill(sql).Config.SQL
		if err := sql.CreateSuperuser(ctx, config.AdminUsername, config.AdminPassword, true); err != nil {
			return fmt.Errorf("failed to create administrative user: %w", err)
		}
	}

	{
		if _, err := logging.LogMessage(progress, "Creating administrative database"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		config := component.GetStill(sql).Config.SQL
		if err := sql.CreateDatabase(ctx, CreateOpts{
			Name:        config.Database,
			AllowExists: true,

			CreateUser: false,
		}); err != nil {
			return fmt.Errorf("failed to create administative database: %w", err)
		}
	}

	// wait for the database to come up
	{
		if _, err := logging.LogMessage(progress, "Waiting for regular database user to connect"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		if err := sql.Wait(ctx); err != nil {
			return fmt.Errorf("failed to wait for database: %w", err)
		}
	}

	// migrate all of the tables!
	if err := logging.LogOperation(func() error {
		for _, table := range sql.dependencies.Tables {
			info := table.TableInfo()
			table := info.Name()

			if _, err := logging.LogMessage(progress, "migrating %q table", table); err != nil {
				return fmt.Errorf("failed to log message: %w", err)
			}

			db, err := sql.connectGorm(ctx)
			if err != nil {
				return fmt.Errorf("failed to connect to table %q: %w", table, err)
			}

			if err := db.AutoMigrate(info.Model); err != nil {
				return fmt.Errorf("failed migration for table %q: %w", table, err)
			}
		}
		return nil
	}, progress, "migrating database tables"); err != nil {
		return fmt.Errorf("failed to migrate database tables: %w", err)
	}
	return nil
}
