package instances

import (
	_ "embed"

	"github.com/tkw1536/goprogram/stream"
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
