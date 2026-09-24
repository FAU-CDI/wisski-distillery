//spellchecker:words extras
package extras

//spellchecker:words context github wisski distillery internal phpx status ingredient embed
import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/stream"

	_ "embed"
)

// Version implements reading the current drupal version.
type Version struct {
	ingredient.Base
	dependencies struct {
		PHP    *php.PHP
		Barrel *barrel.Barrel
	}
}

const versionCode = `return [phpversion(), Drupal::VERSION];`

// Get returns the current versions.
func (v *Version) Get(ctx context.Context, server *phpx.Server) (phpVersion, drupalVersion, wisskiVersion string, err error) {
	var err1, err2 error
	phpVersion, err1 = v.getPHPVersion(ctx, server)
	wisskiVersion, drupalVersion, err2 = v.getVersionsFromLockfile()
	return phpVersion, drupalVersion, wisskiVersion, errorsx.Combine(err1, err2)
}

func (v *Version) getPHPVersion(ctx context.Context, server *phpx.Server) (string, error) {
	// if we already have a server, use that.
	if server != nil {
		var phpVersion string
		if err := v.dependencies.PHP.EvalCode(ctx, server, &phpVersion, "return PHP_VERSION;"); err != nil {
			return "", fmt.Errorf("failed to get php version: %w", err)
		}
		return phpVersion, nil
	}

	// If we don't have a server, it's quicker to directly invoke only the php executable.
	// That avoids any bootstrap.
	var buffer strings.Builder
	if err := v.dependencies.Barrel.BashScript(ctx, stream.IOStream{Stdout: &buffer}, "php", "-r", "echo PHP_VERSION;"); err != nil {
		return "", fmt.Errorf("failed to get php version: %w", err)
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (v *Version) getVersionsFromLockfile() (wisskiVersion string, drupalVersion string, err error) {
	path := filepath.Join(ingredient.GetLiquid(v).FilesystemBase, "data", "data", "project", "composer.lock")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("failed to read composer.lock: %w", err)
	}

	var lock struct {
		Packages []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"packages"`
	}
	if err := json.Unmarshal(data, &lock); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal composer.lock: %w", err)
	}
	for _, pkg := range lock.Packages {
		switch pkg.Name {
		case "drupal/wisski":
			wisskiVersion = pkg.Version
		case "drupal/core":
			drupalVersion = pkg.Version
		}
	}

	return wisskiVersion, drupalVersion, nil
}

func (v *Version) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.PHPVersion, info.DrupalVersion, info.WisskiVersion, err = v.Get(flags.Context, flags.Server)
	if !flags.Quick && err != nil {
		return fmt.Errorf("failed to get versions: %w", err)
	}
	return
}
