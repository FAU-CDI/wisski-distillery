package custom

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

// CustomAssetsPath is the path custom assets are stored at
func (custom *Custom) CustomAssetsPath() string {
	return filepath.Join(custom.Config.DeployRoot, "core", "assets")
}

func (custom *Custom) CustomAssetPath(name string) string {
	return filepath.Join(custom.CustomAssetsPath(), name)
}

func (custom *Custom) BackupName() string { return "custom" }

func (custom *Custom) Backup(context component.StagingContext) error {
	return context.CopyDirectory("", custom.CustomAssetsPath())
}
