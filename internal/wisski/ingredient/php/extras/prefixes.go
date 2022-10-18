package extras

import (
	"bufio"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/tkw1536/goprogram/lib/collection"

	_ "embed"
)

// Prefixes implements reading and writing prefix
type Prefixes struct {
	ingredient.Base

	PHP    *php.PHP
	MStore *mstore.MStore
}

// NoPrefix checks if this WissKI instance is excluded from generating prefixes.
// TODO: Move this to the database!
func (prefixes *Prefixes) NoPrefix() bool {
	return fsx.IsFile(prefixes.Malt.Environment, filepath.Join(prefixes.FilesystemBase, "prefixes.skip"))
}

//go:embed prefixes.php
var listURIPrefixesPHP string

// All returns the prefixes applying to this WissKI
//
// server is an optional server to fetch prefixes from.
// server may be nil.
func (prefixes *Prefixes) All(server *php.Server) ([]string, error) {
	uris, err := prefixes.database(server)
	if err != nil {
		return nil, err
	}

	uris2, err := prefixes.filePrefixes()
	if err != nil {
		return nil, err
	}

	return append(uris, uris2...), nil
}

func (wisski *Prefixes) database(server *php.Server) (prefixes []string, err error) {
	// get all the ugly prefixes
	err = wisski.PHP.ExecScript(server, &prefixes, listURIPrefixesPHP, "list_prefixes")
	if err != nil {
		return nil, err
	}

	// filter out sequential prefixes
	prefixes = collection.NonSequential(prefixes, func(prev, now string) bool {
		return strings.HasPrefix(now, prev)
	})

	// load the list of blocked prefixes
	blocks, err := wisski.blocked()
	if err != nil {
		return nil, err
	}

	// filter out blocked prefixes
	return collection.Filter(prefixes, func(uri string) bool { return !hasAnyPrefix(uri, blocks) }), nil
}

func (prefixes *Prefixes) blocked() ([]string, error) {
	// open the resolver block file
	// TODO: move this to the distillery
	file, err := prefixes.Malt.Environment.Open(prefixes.Malt.Config.SelfResolverBlockFile)
	if err != nil {
		return nil, err
	}

	var lines []string

	// read all the lines that aren't a comment!
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}

	// check if there was an error
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// and done!
	return lines, nil
}

func hasAnyPrefix(candidate string, prefixes []string) bool {
	return collection.Any(
		prefixes,
		func(prefix string) bool {
			return strings.HasPrefix(candidate, prefix)
		},
	)
}

func (wisski *Prefixes) filePrefixes() (prefixes []string, err error) {
	path := filepath.Join(wisski.FilesystemBase, "prefixes")
	if !fsx.IsFile(wisski.Malt.Environment, path) {
		return nil, nil
	}

	file, err := wisski.Malt.Environment.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		prefixes = append(prefixes, line)
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return prefixes, nil
}

// CACHING

var prefix = mstore.For[string]("prefix")

// Prefixes returns the cached prefixes from the given instance
func (wisski *Prefixes) PrefixesCached() (results []string, err error) {
	return prefix.GetAll(wisski.MStore)
}

// UpdatePrefixes updates the cached prefixes of this instance
func (wisski *Prefixes) UpdatePrefixes() error {
	prefixes, err := wisski.All(nil)
	if err != nil {
		return err
	}
	return prefix.SetAll(wisski.MStore, prefixes...)
}
