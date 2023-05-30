package templating

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

// CustomAssetsPath is the path custom assets are stored at
func (tpl *Templating) CustomAssetsPath() string {
	return filepath.Join(tpl.Config.Paths.Root, "core", "assets")
}

func (tpl *Templating) CustomAssetPath(name string) string {
	return filepath.Join(tpl.CustomAssetsPath(), name)
}

func (tpl *Templating) BackupName() string { return "custom" }

func (tpl *Templating) Backup(context *component.StagingContext) error {
	return context.CopyDirectory("", tpl.CustomAssetsPath())
}
