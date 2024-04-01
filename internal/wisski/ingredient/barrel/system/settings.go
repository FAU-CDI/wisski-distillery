package system

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

// RebuildSettings (re-)configures settings.php for the provided running instance
func (smanager *SystemManager) RebuildSettings(ctx context.Context, progress io.Writer) (err error) {
	logging.LogMessage(progress, "Updating TRUSTED_HOST_PATTERNS in settings.php")
	{
		if err := smanager.dependencies.Settings.SetTrustedDomain(ctx, nil, smanager.Domain()); err != nil {
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
