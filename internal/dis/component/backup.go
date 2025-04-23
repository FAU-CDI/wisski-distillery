//spellchecker:words component
package component

//spellchecker:words context path filepath github wisski distillery internal models errors pkglib umaskfree
import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
)

// Backupable represents a component with a Backup method.
type Backupable interface {
	Component

	// BackupName returns a new name to be used as an argument for path.
	BackupName() string

	// Backup backs up this component into the destination path path
	Backup(context *StagingContext) error
}

// Snapshotable represents a component with a Snapshot method.
type Snapshotable interface {
	Component

	// SnapshotNeedsRunning returns if this Snapshotable requires a running instance.
	SnapshotNeedsRunning() bool

	// SnapshotName returns a new name to be used as an argument for path.
	SnapshotName() string

	// Snapshot snapshots a part of the instance
	Snapshot(wisski models.Instance, context *StagingContext) error
}

// NewStagingContext returns a new [StagingContext].
func NewStagingContext(ctx context.Context, progress io.Writer, path string, manifest chan<- string) *StagingContext {
	return &StagingContext{
		ctx:      ctx,
		progress: progress,
		path:     path,
		manifest: manifest,
	}
}

// StagingContext is a context used for [Backupable] and [Snapshotable].
//
//nolint:containedctx
type StagingContext struct {
	ctx      context.Context
	progress io.Writer     // writer to direct progress to
	path     string        // path to send files to
	manifest chan<- string // channel the manifest is sent to
}

func (bc *StagingContext) sendPath(path string) {
	// ensure path is absolute!
	if !filepath.IsAbs(path) {
		var err error
		path, err = bc.resolve(path)
		if err != nil {
			fmt.Fprintf(bc.progress, "path resolve error: %s", err)
			return
		}
	}

	// use the relative path for logging
	rel, err := bc.relativize(path)
	if err == nil {
		_, _ = io.WriteString(bc.progress, rel+"\n")
	}

	// send the absolute path
	bc.manifest <- path
}

// Progress returns a writer to write progress information to.
func (bc *StagingContext) Progress() io.Writer {
	return bc.progress
}

var (
	errResolveAbsolute  = errors.New("resolve: path must be relative")
	errRelativeRelative = errors.New("relativize: path already relative")
)

func (bc *StagingContext) resolve(path string) (dest string, err error) {
	if path == "" {
		return bc.path, nil
	}
	if filepath.IsAbs(path) {
		return "", errResolveAbsolute
	}
	return filepath.Join(bc.path, path), nil
}

func (bc *StagingContext) relativize(path string) (dest string, err error) {
	if !filepath.IsAbs(path) {
		return "", errRelativeRelative
	}

	return filepath.Rel(bc.path, path)
}

// AddDirectory creates a new directory inside the destination.
// Passing the empty path creates the destination as a directory.
//
// It then allows op to fill the file.
func (sc *StagingContext) AddDirectory(path string, op func(context.Context) error) error {
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

// CopyFile copies a file from src to dst.
func (sc *StagingContext) CopyFile(dst, src string) error {
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

// CopyDirectory copies a directory from src to dst.
func (sc *StagingContext) CopyDirectory(dst, src string) error {
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

// AddFile creates a new file at the provided path inside the destination.
// Passing the empty path creates the destination as a file.
//
// It then allows op to write to the file.
//
// The op function must not retain file.
// The underlying file does not need to be closed.
// AddFile will not return before op has returned.
func (sc *StagingContext) AddFile(path string, op func(ctx context.Context, file io.Writer) error) error {
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

func (sc *StagingContext) ctxdone() (err error, done bool) {
	err = sc.ctx.Err()
	done = (err != nil)
	return
}
