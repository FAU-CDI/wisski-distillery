package control

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/tkw1536/goprogram/stream"
)

func (*Control) BackupName() string {
	return "config"
}

// Backup backups all control plane configuration files into dest
func (control *Control) Backup(io stream.IOStream, dest string) error {
	// create the destination directory, TODO: outsource this
	if err := os.Mkdir(dest, fs.ModeDir); err != nil {
		return err
	}

	files := control.backupFiles()
	for _, src := range files {
		dst := filepath.Join(dest, filepath.Base(src)) // destination path

		// if the src file does not exist, don't copy it!
		if !fsx.IsFile(src) { // TODO: log this somewhere
			continue
		}

		if err := fsx.CopyFile(dst, src); err != nil {
			return err
		}
	}

	return nil
}

// backupfiles lists the files to be backed up.
func (control *Control) backupFiles() []string {
	return []string{
		control.Config.ConfigPath,
		control.Config.ExecutablePath(),
		control.Config.SelfOverridesFile,
		control.Config.GlobalAuthorizedKeysFile,
	}
}
