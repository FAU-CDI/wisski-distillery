package instances

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var errPathbuildersExecFailed = errors.New("ExportPathbuilders: Failed to call export_pathbuilder")

// ExportPathbuilders writes pathbuilders into the directory dest
func (wisski *WissKI) ExportPathbuilders(dest string) error {
	// export all the pathbuilders into the buffer
	var buffer bytes.Buffer
	wu := stream.NewIOStream(&buffer, nil, nil, 0)
	code, err := wisski.Barrel().Exec(wu, "barrel", "/bin/bash", "/user_shell.sh", "-c", "drush php:script /wisskiutils/export_pathbuilder.php")
	if err != nil || code != 0 {
		return errPathbuildersExecFailed
	}

	// decode them as a json array
	var pathbuilders map[string]string
	if err := json.NewDecoder(&buffer).Decode(&pathbuilders); err != nil {
		return err
	}

	// sort the names of the pathbuilders
	names := maps.Keys(pathbuilders)
	slices.Sort(names)

	// write each into a file!
	for _, name := range names {
		pbxml := []byte(pathbuilders[name])
		name := filepath.Join(dest, fmt.Sprintf("%s.xml", name))
		if err := os.WriteFile(name, pbxml, fs.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
