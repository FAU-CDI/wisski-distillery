package backup

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/tkw1536/goprogram/stream"
)

// context implements [components.BackupContext]
type context struct {
	env   environment.Environment
	io    stream.IOStream
	dst   string      // destination directory
	files chan string // files channel
}

func (bc *context) sendPath(path string) {

	// resolve the path, or bail out!
	// TODO: Use the relative path here!
	dst, err := bc.resolve(path)
	if err != nil {
		return
	}

	bc.io.Println(dst)
	bc.files <- dst
}

func (bc *context) IO() stream.IOStream {
	return bc.io
}

var errResolveAbsolute = errors.New("resolve: path must be relative")

func (bc *context) resolve(path string) (dest string, err error) {
	if path == "" {
		return bc.dst, nil
	}
	if filepath.IsAbs(path) {
		return "", errResolveAbsolute
	}
	return filepath.Join(bc.dst, path), nil
}

func (bc *context) AddDirectory(path string, op func() error) error {
	// resolve the path!
	dst, err := bc.resolve(path)
	if err != nil {
		return err
	}

	// run the make directory
	if err := bc.env.Mkdir(dst, environment.DefaultDirPerm); err != nil {
		return err
	}

	// tell the files that we are creating it!
	bc.sendPath(path)

	// and run the files!
	// TODO: Add to manifest of some sort
	return op()
}

func (bc *context) CopyFile(dst, src string) error {
	dstPath, err := bc.resolve(dst)
	if err != nil {
		return err
	}
	bc.sendPath(dst)
	return fsx.CopyFile(bc.env, dstPath, src)
}

func (bc *context) AddFile(path string, op func(file io.Writer) error) error {
	// resolve the path!
	dst, err := bc.resolve(path)
	if err != nil {
		return err
	}

	// create the file
	file, err := bc.env.Create(dst, environment.DefaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	// tell them that we are creating it!
	bc.sendPath(path)

	// and do whatever they wanted to do
	return op(file)
}
