package exporter

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"

	"github.com/FAU-CDI/wisski-distillery/internal/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/password"
)

// Exporter manages snapshots and backups
type Exporter struct {
	component.ComponentBase

	SQL            *sql.SQL
	Instances      *instances.Instances
	ExporterLogger *logger.Logger

	Snapshotable []component.Snapshotable
	Backupable   []component.Backupable
}

func (Exporter) Name() string { return "snapshots" }

// Path returns the path that contains all snapshot related data.
func (dis *Exporter) Path() string {
	return filepath.Join(dis.Config.DeployRoot, "snapshots")
}

// StagingPath returns the path to the directory containing a temporary staging area for snapshots.
// Use NewSnapshotStagingDir to generate a new staging area.
func (dis *Exporter) StagingPath() string {
	return filepath.Join(dis.Path(), "staging")
}

// ArchivePath returns the path to the directory containing all exported archives.
// Use NewSnapshotArchivePath to generate a path to a new archive in this directory.
func (dis *Exporter) ArchivePath() string {
	return filepath.Join(dis.Path(), "archives")
}

// NewArchivePath returns the path to a new archive with the provided prefix.
// The path is guaranteed to not exist.
func (dis *Exporter) NewArchivePath(prefix string) (path string) {
	// TODO: Consider moving these into a subdirectory with the provided prefix.
	for path == "" || fsx.Exists(dis.Environment, path) {
		name := dis.newSnapshotName(prefix) + ".tar.gz"
		path = filepath.Join(dis.ArchivePath(), name)
	}
	return
}

// newSnapshot name returns a new basename for a snapshot with the provided prefix.
// The name is guaranteed to be unique within this process.
func (*Exporter) newSnapshotName(prefix string) string {
	suffix, _ := password.Password(10) // silently ignore any errors!
	if prefix == "" {
		prefix = "backup"
	} else {
		prefix = "snapshot-" + prefix
	}
	return fmt.Sprintf("%s-%d-%s", prefix, time.Now().Unix(), suffix)
}

// NewStagingDir returns the path to a new snapshot directory.
// The directory is guaranteed to have been freshly created.
func (dis *Exporter) NewStagingDir(prefix string) (path string, err error) {
	for path == "" || environment.IsExist(err) {
		path = filepath.Join(dis.StagingPath(), dis.newSnapshotName(prefix))
		err = dis.Core.Environment.Mkdir(path, environment.DefaultFilePerm)
	}
	if err != nil {
		path = ""
	}
	return
}
