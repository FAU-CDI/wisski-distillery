//spellchecker:words templating
package templating

//spellchecker:words path filepath github wisski distillery internal component
import (
	"fmt"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

// CustomAssetsPath is the path custom assets are stored at.
func (tpl *Templating) CustomAssetsPath() string {
	return filepath.Join(component.GetStill(tpl).Config.Paths.Root, "core", "assets")
}

func (tpl *Templating) CustomAssetPath(name string) string {
	return filepath.Join(tpl.CustomAssetsPath(), name)
}

func (tpl *Templating) BackupName() string { return "custom" }

func (tpl *Templating) Backup(context *component.StagingContext) error {
	if err := context.CopyDirectory("", tpl.CustomAssetsPath()); err != nil {
		return fmt.Errorf("failed to copy custom assets: %w", err)
	}
	return nil
}
