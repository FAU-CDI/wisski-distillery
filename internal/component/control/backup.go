package control

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
)

func (*Control) BackupName() string {
	return "config"
}

// Backup backups all control plane configuration files into dest
func (control *Control) Backup(context component.StagingContext) error {
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
func (control *Control) backupFiles() []string {
	return []string{
		control.Config.ConfigPath,
		control.Config.ExecutablePath(),
		control.Config.SelfOverridesFile,
		control.Config.SelfResolverBlockFile,
		control.Config.GlobalAuthorizedKeysFile,
	}
}
