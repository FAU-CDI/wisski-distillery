package instances

import (
	_ "embed"

	"golang.org/x/exp/slices"
)

//go:embed php/export_pathbuilder.php
var exportPathbuilderPHP string

// Pathbuilders returns the ids of all pathbuilders in consistent order.
//
// server is the server to fetch the pathbuilders from, any may be nil.
func (wisski *WissKI) Pathbuilders(server *PHPServer) (ids []string, err error) {
	err = wisski.ExecPHPScript(server, &ids, exportPathbuilderPHP, "all_list")
	slices.Sort(ids)
	return
}

// Pathbuilder returns a single pathbuilder as xml.
// If it does not exist, it returns the empty string and nil error.
//
// server is the server to fetch the pathbuilders from, any may be nil.
func (wisski *WissKI) Pathbuilder(server *PHPServer, id string) (xml string, err error) {
	err = wisski.ExecPHPScript(server, &xml, exportPathbuilderPHP, "one_xml", id)
	return
}

// AllPathbuilders returns all pathbuilders serialized as xml
//
// server is the server to fetch the pathbuilders from, any may be nil.
func (wisski *WissKI) AllPathbuilders(server *PHPServer) (pathbuilders map[string]string, err error) {
	err = wisski.ExecPHPScript(server, &pathbuilders, exportPathbuilderPHP, "all_xml")
	return
}
