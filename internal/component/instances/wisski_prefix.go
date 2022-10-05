package instances

import (
	"bufio"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/slicesx"
	"github.com/tkw1536/goprogram/stream"

	_ "embed"
)

// NoPrefix checks if this WissKI instance is excluded from generating prefixes.
// TODO: Move this to the database!
func (wisski *WissKI) NoPrefix() bool {
	return fsx.IsFile(wisski.instances.Environment, filepath.Join(wisski.FilesystemBase, "prefixes.skip"))
}

//go:embed php/list_uri_prefixes.php
var listURIPrefixesPHP string

// Prefixes returns the prefixes applying to this WissKI
func (wisski *WissKI) Prefixes() ([]string, error) {
	prefixes, err := wisski.dbPrefixes()
	if err != nil {
		return nil, err
	}

	prefixes2, err := wisski.filePrefixes()
	if err != nil {
		return nil, err
	}

	return append(prefixes, prefixes2...), nil
}

func (wisski *WissKI) dbPrefixes() (prefixes []string, err error) {
	// get all the ugly prefixes
	err = wisski.ExecPHPScript(stream.FromDebug(), &prefixes, listURIPrefixesPHP, "list_prefixes")
	if err != nil {
		return nil, err
	}

	// filter out sequential prefixes
	prefixes = slicesx.NonSequential(prefixes, func(prev, now string) bool {
		return strings.HasPrefix(now, prev)
	})

	// load the list of blocked prefixes
	blocks, err := wisski.instances.blockedPrefixes()
	if err != nil {
		return nil, err
	}

	// filter out blocked prefixes
	return slicesx.Filter(prefixes, func(uri string) bool { return !hasAnyPrefix(uri, blocks) }), nil
}

func (instances *Instances) blockedPrefixes() ([]string, error) {
	// open the resolver block file
	file, err := instances.Environment.Open(instances.Config.SelfResolverBlockFile)
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
	return slicesx.Any(
		prefixes,
		func(prefix string) bool {
			return strings.HasPrefix(candidate, prefix)
		},
	)
}

func (wisski *WissKI) filePrefixes() (prefixes []string, err error) {
	path := filepath.Join(wisski.FilesystemBase, "prefixes")
	if !fsx.IsFile(wisski.instances.Environment, path) {
		return nil, nil
	}

	file, err := wisski.instances.Environment.Open(path)
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

var PrefixConfigKey MetaKey = "prefix"

// Prefixes returns the cached prefixes from the given instance
func (wisski *WissKI) PrefixesCached() (results []string, err error) {
	err = wisski.Metadata().GetAll(PrefixConfigKey, func(index, total int) any {
		if results == nil {
			results = make([]string, total)
		}
		return &results[index]
	})
	return
}

// UpdatePrefixes updates the cached prefixes of this instance
func (wisski *WissKI) UpdatePrefixes() error {
	prefixes, err := wisski.Prefixes()
	if err != nil {
		return err
	}

	return wisski.Metadata().SetAll(PrefixConfigKey, slicesx.AsAny(prefixes)...)
}
