package exporter

import (
	"context"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

// Config implements backing up configuration
type Config struct {
	component.Base
}

var (
	_ = (component.Backupable)((*Config)(nil))
)

func (*Config) BackupName() string {
	return "config"
}

func (control *Config) Backup(scontext *component.StagingContext) error {
	files := control.backupFiles()

	return scontext.AddDirectory("", func(ctx context.Context) error {
		for _, src := range files {
			name := filepath.Base(src)
			if err := scontext.CopyFile(name, src); err != nil {
				return err
			}
		}
		return nil
	})
}

// backupfiles lists the files to be backed up.
func (control *Config) backupFiles() []string {
	config := component.GetStill(control).Config
	return []string{
		config.ConfigPath,
		config.Paths.ExecutablePath(),
		config.Paths.OverridesJSON,
		config.Paths.ResolverBlocks,
	}
}
