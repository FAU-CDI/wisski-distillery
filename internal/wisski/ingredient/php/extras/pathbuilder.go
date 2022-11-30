package extras

import (
	"context"
	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"golang.org/x/exp/slices"
)

type Pathbuilder struct {
	ingredient.Base

	PHP *php.PHP
}

var (
	_ ingredient.WissKIFetcher = (*Pathbuilder)(nil)
)

//go:embed pathbuilder.php
var pathbuilderPHP string

// All returns the ids of all pathbuilders in consistent order.
//
// server is the server to fetch the pathbuilders from, any may be nil.
func (pathbuilder *Pathbuilder) All(ctx context.Context, server *phpx.Server) (ids []string, err error) {
	err = pathbuilder.PHP.ExecScript(ctx, server, &ids, pathbuilderPHP, "all_list")
	slices.Sort(ids)
	return
}

// Get returns a single pathbuilder as xml.
// If it does not exist, it returns the empty string and nil error.
//
// server is the server to fetch the pathbuilders from, any may be nil.
func (pathbuilder *Pathbuilder) Get(ctx context.Context, server *phpx.Server, id string) (xml string, err error) {
	err = pathbuilder.PHP.ExecScript(ctx, server, &xml, pathbuilderPHP, "one_xml", id)
	return
}

// GetAll returns all pathbuilders serialized as xml
//
// server is the server to fetch the pathbuilders from, any may be nil.
func (pathbuilder *Pathbuilder) GetAll(ctx context.Context, server *phpx.Server) (pathbuilders map[string]string, err error) {
	err = pathbuilder.PHP.ExecScript(ctx, server, &pathbuilders, pathbuilderPHP, "all_xml")
	return
}

func (pathbuilder *Pathbuilder) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.Pathbuilders, _ = pathbuilder.GetAll(flags.Context, flags.Server)
	return
}
