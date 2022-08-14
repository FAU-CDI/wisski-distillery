package env

import "path/filepath"

func (dis Distillery) BackupDir() string {
	return filepath.Join(dis.Config.DeployRoot, "backups")
}

func (dis Distillery) RuntimeDir() string {
	return filepath.Join(dis.Config.DeployRoot, "runtime")
}

func (dis Distillery) RuntimeUtilsDir() string {
	return filepath.Join(dis.Config.DeployRoot, "runtime", "utils")
}

func (dis Distillery) InprogressBackupPath() string {
	return filepath.Join(dis.BackupDir(), "inprogress")
}

func (dis Distillery) FinalBackupPath() string {
	return filepath.Join(dis.BackupDir(), "final")
}
