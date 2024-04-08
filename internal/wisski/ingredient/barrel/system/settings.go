package system

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

// BuildSettings sets up global settings.php configuration settings.php for the provided running instance
// This doesn't need to be called manually.
func (smanager *SystemManager) BuildSettings(ctx context.Context, progress io.Writer) (err error) {
	logging.LogMessage(progress, "Updating TRUSTED_HOST_PATTERNS in settings.php")
	{
		if err := smanager.dependencies.Settings.SetTrustedDomain(ctx, nil, ingredient.GetLiquid(smanager).Domain()); err != nil {
			return err
		}
	}

	logging.LogMessage(progress, "Adding distillery settings to settings.php")
	{
		if err := smanager.dependencies.Settings.InstallDistillerySettings(ctx, nil); err != nil {
			return err
		}
	}

	return nil
}
