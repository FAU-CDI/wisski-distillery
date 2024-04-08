package exporter

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/passwordx"
	"github.com/tkw1536/pkglib/fsx"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
	"github.com/tkw1536/pkglib/password"
)

// Exporter manages snapshots and backups
type Exporter struct {
	component.Base
	dependencies struct {
		SQL            *sql.SQL
		Instances      *instances.Instances
		ExporterLogger *logger.Logger

		Snapshotable []component.Snapshotable
		Backupable   []component.Backupable
	}
}

// Path returns the path that contains all snapshot related data.
func (dis *Exporter) Path() string {
	return filepath.Join(component.GetStill(dis).Config.Paths.Root, "snapshots")
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
	for {
		name := dis.newSnapshotName(prefix) + ".tar.gz"
		path = filepath.Join(dis.ArchivePath(), name)

		// Only return if the path does not exist
		// FIXME: ignoring the error here
		exists, _ := fsx.Exists(path)
		if !exists {
			return path
		}
	}
}

// newSnapshot name returns a new basename for a snapshot with the provided prefix.
// The name is guaranteed to be unique within this process.
func (*Exporter) newSnapshotName(prefix string) string {
	suffix, _ := password.Generate(rand.Reader, 10, passwordx.Identifier) // silently ignore any errors!
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
	for path == "" || errors.Is(err, fs.ErrExist) {
		path = filepath.Join(dis.StagingPath(), dis.newSnapshotName(prefix))
		err = umaskfree.Mkdir(path, umaskfree.DefaultFilePerm)
	}
	if err != nil {
		path = ""
	}
	return
}
