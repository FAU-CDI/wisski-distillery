package system

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
)

// SystemManager applies a specific system configuration
type SystemManager struct {
	ingredient.Base
	Dependencies struct {
		Barrel      *barrel.Barrel
		Bookkeeping *bookkeeping.Bookkeeping
		Settings    *extras.Settings
	}
}

// Apply applies a specific system version to this barrel.
// If start is true, also starts the container.
func (smanager *SystemManager) Apply(ctx context.Context, progress io.Writer, system models.System, start bool) (err error) {
	// setup the new docker image
	smanager.Instance.System = system

	// save in bookkeeping
	if err := smanager.Dependencies.Bookkeeping.Save(ctx); err != nil {
		return err
	}

	// TODO: Apply Content-Security-Policy!

	// and rebuild
	return smanager.Dependencies.Barrel.Build(ctx, progress, start)
}
