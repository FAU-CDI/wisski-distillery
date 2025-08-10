//spellchecker:words extras
package extras

//spellchecker:words context github wisski distillery internal phpx status ingredient embed
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"

	_ "embed"
)

// Version implements reading the current drupal version.
type Version struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

const versionCode = `return [phpversion(), Drupal::VERSION];`

// Get returns the currently active theme.
func (v *Version) Get(ctx context.Context, server *phpx.Server) (phpVersion, drupalVersion string, err error) {
	var versions [2]string
	err = v.dependencies.PHP.EvalCode(
		ctx, server, &versions, versionCode,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to get versions: %w", err)
	}
	phpVersion = versions[0]
	drupalVersion = versions[1]
	return
}

func (v *Version) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.PHPVersion, info.DrupalVersion, err = v.Get(flags.Context, flags.Server)
	if err != nil {
		return fmt.Errorf("failed to get versions: %w", err)
	}
	return
}
