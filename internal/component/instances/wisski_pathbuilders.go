package instances

import (
	"fmt"
	"path/filepath"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

//go:embed php/export_pathbuilder.php
var exportPathbuilderPHP string

// Pathbuilders returns the ids of all pathbuilders in consistent order.
func (wisski *WissKI) Pathbuilders() (ids []string, err error) {
	err = wisski.ExecPHPScript(stream.FromDebug(), &ids, exportPathbuilderPHP, "all_list")
	slices.Sort(ids)
	return
}

// Pathbuilder returns a single pathbuilder as xml.
// If it does not exist, it returns the empty string and nil error.
func (wisski *WissKI) Pathbuilder(id string) (xml string, err error) {
	err = wisski.ExecPHPScript(stream.FromDebug(), &xml, exportPathbuilderPHP, "one_xml", id)
	return
}

// AllPathbuilders returns all pathbuilders serialized as xml
func (wisski *WissKI) AllPathbuilders() (pathbuilders map[string]string, err error) {
	err = wisski.ExecPHPScript(stream.FromDebug(), &pathbuilders, exportPathbuilderPHP, "all_xml")
	return
}

// ExportPathbuilders writes pathbuilders into the directory dest
func (wisski *WissKI) ExportPathbuilders(dest string) error {
	pathbuilders, err := wisski.AllPathbuilders()
	if err != nil {
		return err
	}

	// sort the names of the pathbuilders
	names := maps.Keys(pathbuilders)
	slices.Sort(names)

	// write each into a file!
	for _, name := range names {
		pbxml := []byte(pathbuilders[name])
		name := filepath.Join(dest, fmt.Sprintf("%s.xml", name))
		if err := environment.WriteFile(wisski.instances.Core.Environment, name, pbxml, environment.DefaultFilePerm); err != nil {
			return err
		}
	}

	return nil
}
