//spellchecker:words system
package system

//spellchecker:words context github wisski distillery internal models ingredient barrel bookkeeping extras
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
)

// SystemManager applies a specific system configuration.
type SystemManager struct {
	ingredient.Base
	dependencies struct {
		Barrel      *barrel.Barrel
		Bookkeeping *bookkeeping.Bookkeeping
		Settings    *extras.Settings
	}
}

// Apply applies the given system configuration to this instance and (re-)starts the system.
func (smanager *SystemManager) Apply(ctx context.Context, progress io.Writer, system models.System) (err error) {
	if err := smanager.apply(ctx, progress, system, true); err != nil {
		return err
	}

	if err := smanager.BuildSettings(ctx, progress); err != nil {
		return err
	}

	return nil
}

// ApplyInitial builds the base image, but does not start it.
func (smanager *SystemManager) ApplyInitial(ctx context.Context, progress io.Writer, system models.System) error {
	return smanager.apply(ctx, progress, system, false)
}

// start inidicates if the image should be started afterwards.
func (smanager *SystemManager) apply(ctx context.Context, progress io.Writer, system models.System, start bool) error {
	// store the new system configuration
	ingredient.GetLiquid(smanager).System = system
	if err := smanager.dependencies.Bookkeeping.Save(ctx); err != nil {
		return fmt.Errorf("failed to save bookkeeping: %w", err)
	}

	// build and start the barrel
	if err := smanager.dependencies.Barrel.Build(ctx, progress, start); err != nil {
		return fmt.Errorf("faield to build barrel: %w", err)
	}
	return nil
}
