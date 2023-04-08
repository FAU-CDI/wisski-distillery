package component

import (
	"context"
	"io"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/pkg/errors"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
)

// Backupable represents a component with a Backup method
type Backupable interface {
	Component

	// BackupName returns a new name to be used as an argument for path.
	BackupName() string

	// Backup backs up this component into the destination path path
	Backup(context StagingContext) error
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

// StagingContext represents a context for [Backupable] and [Snapshotable]
type StagingContext interface {
	// Progress returns a writer to write progress information to.
	Progress() io.Writer

	// Name creates a new directory inside the destination.
	// Passing the empty path creates the destination as a directory.
	//
	// It then allows op to fill the file.
	AddDirectory(path string, op func(context.Context) error) error

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
	AddFile(path string, op func(ctx context.Context, file io.Writer) error) error
}

// NewStagingContext returns a new [StagingContext]
func NewStagingContext(ctx context.Context, progress io.Writer, path string, manifest chan<- string) StagingContext {
	return &stagingContext{
		ctx:      ctx,
		progress: progress,
		path:     path,
		manifest: manifest,
	}
}

// stagingContext implements [components.StagingContext]
type stagingContext struct {
	ctx      context.Context
	progress io.Writer     // writer to direct progress to
	path     string        // path to send files to
	manifest chan<- string // channel the manifest is sent to
}

func (bc *stagingContext) sendPath(path string) {
	// resolve the path, or bail out!
	// TODO: Use the relative path here!
	dst, err := bc.resolve(path)
	if err != nil {
		return
	}

	io.WriteString(bc.progress, dst+"\n")
	bc.manifest <- dst
}

func (bc *stagingContext) Progress() io.Writer {
	return bc.progress
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

func (sc *stagingContext) AddDirectory(path string, op func(context.Context) error) error {
	// check if we are already done
	if err, ok := sc.ctxdone(); ok {
		return err
	}

	// resolve the path!
	dst, err := sc.resolve(path)
	if err != nil {
		return err
	}

	// run the make directory
	if err := umaskfree.Mkdir(dst, umaskfree.DefaultDirPerm); err != nil {
		return err
	}

	// tell the files that we are creating it!
	sc.sendPath(path)

	// and run the files!
	return op(sc.ctx)
}

func (sc *stagingContext) CopyFile(dst, src string) error {
	if err, ok := sc.ctxdone(); ok {
		return err
	}

	dstPath, err := sc.resolve(dst)
	if err != nil {
		return err
	}
	sc.sendPath(dst)
	return umaskfree.CopyFile(sc.ctx, dstPath, src)
}

func (sc *stagingContext) CopyDirectory(dst, src string) error {
	if err, ok := sc.ctxdone(); ok {
		return err
	}

	dstPath, err := sc.resolve(dst)
	if err != nil {
		return err
	}

	return umaskfree.CopyDirectory(sc.ctx, dstPath, src, func(dst, src string) {
		sc.sendPath(dst)
	})
}

func (sc *stagingContext) AddFile(path string, op func(ctx context.Context, file io.Writer) error) error {
	// check if we're already done
	if err, ok := sc.ctxdone(); ok {
		return err
	}

	// resolve the path!
	dst, err := sc.resolve(path)
	if err != nil {
		return err
	}

	// create the file
	file, err := umaskfree.Create(dst, umaskfree.DefaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// tell them that we are creating it!
	sc.sendPath(path)

	// and do whatever they wanted to do
	return op(sc.ctx, file)
}

func (sc *stagingContext) ctxdone() (err error, done bool) {
	err = sc.ctx.Err()
	done = (err != nil)
	return
}
