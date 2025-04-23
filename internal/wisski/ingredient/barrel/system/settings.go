//spellchecker:words system
package system

//spellchecker:words context github wisski distillery internal ingredient logging
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

// BuildSettings sets up global settings.php configuration settings.php for the provided running instance
// This doesn't need to be called manually.
func (smanager *SystemManager) BuildSettings(ctx context.Context, progress io.Writer) (err error) {
	if _, err := logging.LogMessage(progress, "Updating TRUSTED_HOST_PATTERNS in settings.php"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	{
		if err := smanager.dependencies.Settings.SetTrustedDomain(ctx, nil, ingredient.GetLiquid(smanager).Domain()); err != nil {
			return fmt.Errorf("failed to set trusted domain: %w", err)
		}
	}

	if _, err := logging.LogMessage(progress, "Adding distillery settings to settings.php"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	{
		if err := smanager.dependencies.Settings.InstallDistillerySettings(ctx, nil); err != nil {
			return fmt.Errorf("failed to install distillery settings: %w", err)
		}
	}

	return nil
}
