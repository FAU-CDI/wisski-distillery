package wisski

import (
	"bufio"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/tkw1536/goprogram/lib/collection"

	_ "embed"
)

// NoPrefix checks if this WissKI instance is excluded from generating prefixes.
// TODO: Move this to the database!
func (wisski *WissKI) NoPrefix() bool {
	return fsx.IsFile(wisski.Core.Environment, filepath.Join(wisski.FilesystemBase, "prefixes.skip"))
}

//go:embed php/list_uri_prefixes.php
var listURIPrefixesPHP string

// Prefixes returns the prefixes applying to this WissKI
//
// server is an optional server to fetch prefixes from.
// server may be nil.
func (wisski *WissKI) Prefixes(server *PHPServer) ([]string, error) {
	prefixes, err := wisski.dbPrefixes(server)
	if err != nil {
		return nil, err
	}

	prefixes2, err := wisski.filePrefixes()
	if err != nil {
		return nil, err
	}

	return append(prefixes, prefixes2...), nil
}

func (wisski *WissKI) dbPrefixes(server *PHPServer) (prefixes []string, err error) {
	// get all the ugly prefixes
	err = wisski.ExecPHPScript(server, &prefixes, listURIPrefixesPHP, "list_prefixes")
	if err != nil {
		return nil, err
	}

	// filter out sequential prefixes
	prefixes = collection.NonSequential(prefixes, func(prev, now string) bool {
		return strings.HasPrefix(now, prev)
	})

	// load the list of blocked prefixes
	blocks, err := wisski.blockedPrefixes()
	if err != nil {
		return nil, err
	}

	// filter out blocked prefixes
	return collection.Filter(prefixes, func(uri string) bool { return !hasAnyPrefix(uri, blocks) }), nil
}

func (wisski *WissKI) blockedPrefixes() ([]string, error) {
	// open the resolver block file
	file, err := wisski.Core.Environment.Open(wisski.Core.Config.SelfResolverBlockFile)
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

func (wisski *WissKI) filePrefixes() (prefixes []string, err error) {
	path := filepath.Join(wisski.FilesystemBase, "prefixes")
	if !fsx.IsFile(wisski.Core.Environment, path) {
		return nil, nil
	}

	file, err := wisski.Core.Environment.Open(path)
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

var prefix = meta.StorageFor[string]("prefix")

// Prefixes returns the cached prefixes from the given instance
func (wisski *WissKI) PrefixesCached() (results []string, err error) {
	return prefix(wisski.storage()).GetAll()
}

// UpdatePrefixes updates the cached prefixes of this instance
func (wisski *WissKI) UpdatePrefixes() error {
	prefixes, err := wisski.Prefixes(nil)
	if err != nil {
		return err
	}
	return prefix(wisski.storage()).SetAll(prefixes...)
}
