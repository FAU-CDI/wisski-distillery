package extras

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
)

// Config implements backing up configuration
type Config struct {
	component.ComponentBase
}

func (Config) Name() string { return "extra-config" }

func (*Config) BackupName() string {
	return "config"
}

func (control *Config) Backup(context component.StagingContext) error {
	files := control.backupFiles()

	return context.AddDirectory("", func() error {
		for _, src := range files {
			name := filepath.Base(src)
			if err := context.CopyFile(name, src); err != nil {
				return err
			}
		}
		return nil
	})
}

// backupfiles lists the files to be backed up.
func (control *Config) backupFiles() []string {
	return []string{
		control.Config.ConfigPath,
		control.Config.ExecutablePath(),
		control.Config.SelfOverridesFile,
		control.Config.SelfResolverBlockFile,
		control.Config.GlobalAuthorizedKeysFile,
	}
}
