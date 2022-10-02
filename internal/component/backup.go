package component

import (
	"io"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/stream"
)

// Backupable represents a component with a Backup method
type Backupable interface {
	Component

	// BackupName returns a new name to be used as an argument for path.
	BackupName() string

	// Backup backs up this component into the destination path path
	Backup(context StagingContext) error
}

// StagingContext represents a context for [Backupable] and [Snapshotable]
type StagingContext interface {
	// IO returns the input output stream belonging to this backup file
	IO() stream.IOStream

	// Name creates a new directory inside the destination.
	// Passing the empty path creates the destination as a directory.
	//
	// It then allows op to fill the file.
	AddDirectory(path string, op func() error) error

	// CopyFile copies a file from src to dst.
	CopyFile(dst, src string) error

	// CopyDirectory copies a directory from src to dst.
	CopyDirectory(dst, src string) error

	// AddFile creates a new file at the provided path inside the destination.
	// Passing the empty path creates the destination as a file.
	//
	// It then allows op to write to the file.
	//
	// The op function must not retain file.
	// The underlying file does not need to be closed.
	// AddFile will not return before op has returned.
	AddFile(path string, op func(file io.Writer) error) error
}

// Snapshotable represents a component with a Snapshot method.
type Snapshotable interface {
	Component

	// SnapshotNeedsRunning returns if this Snapshotable requires a running instance.
	SnapshotNeedsRunning() bool

	// SnapshotName returns a new name to be used as an argument for path.
	SnapshotName() string

	// Snapshot snapshots a part of the instance
	Snapshot(wisski models.Instance, context StagingContext) error
}

// NewStagingContext returns a new [StagingContext]
func NewStagingContext(env environment.Environment, io stream.IOStream, path string, manifest chan<- string) StagingContext {
	return &stagingContext{
		env:      env,
		io:       io,
		path:     path,
		manifest: manifest,
	}
}

// stagingContext implements [components.StagingContext]
type stagingContext struct {
	env      environment.Environment // environment
	io       stream.IOStream         // context the files are sent to
	path     string                  // path to send files to
	manifest chan<- string           // channel the manifest is sent to
}

func (bc *stagingContext) sendPath(path string) {
	// resolve the path, or bail out!
	// TODO: Use the relative path here!
	dst, err := bc.resolve(path)
	if err != nil {
		return
	}

	bc.io.Println(dst)
	bc.manifest <- dst
}

func (bc *stagingContext) IO() stream.IOStream {
	return bc.io
}

var errResolveAbsolute = errors.New("resolve: path must be relative")

func (bc *stagingContext) resolve(path string) (dest string, err error) {
	if path == "" {
		return bc.path, nil
	}
	if filepath.IsAbs(path) {
		return "", errResolveAbsolute
	}
	return filepath.Join(bc.path, path), nil
}

func (sc *stagingContext) AddDirectory(path string, op func() error) error {
	// resolve the path!
	dst, err := sc.resolve(path)
	if err != nil {
		return err
	}

	// run the make directory
	if err := sc.env.Mkdir(dst, environment.DefaultDirPerm); err != nil {
		return err
	}

	// tell the files that we are creating it!
	sc.sendPath(path)

	// and run the files!
	return op()
}

func (sc *stagingContext) CopyFile(dst, src string) error {
	dstPath, err := sc.resolve(dst)
	if err != nil {
		return err
	}
	sc.sendPath(dst)
	return fsx.CopyFile(sc.env, dstPath, src)
}

func (sc *stagingContext) CopyDirectory(dst, src string) error {
	dstPath, err := sc.resolve(dst)
	if err != nil {
		return err
	}

	return fsx.CopyDirectory(sc.env, dstPath, src, func(dst, src string) {
		sc.sendPath(dst)
	})
}

func (sc *stagingContext) AddFile(path string, op func(file io.Writer) error) error {
	// resolve the path!
	dst, err := sc.resolve(path)
	if err != nil {
		return err
	}

	// create the file
	file, err := sc.env.Create(dst, environment.DefaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// tell them that we are creating it!
	sc.sendPath(path)

	// and do whatever they wanted to do
	return op(file)
}
