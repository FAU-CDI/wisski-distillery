package extras

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"github.com/tkw1536/pkglib/collection"
	"github.com/tkw1536/pkglib/fsx"
	"golang.org/x/exp/slices"

	_ "embed"
)

// Prefixes implements reading and writing prefix
type Prefixes struct {
	ingredient.Base
	dependencies struct {
		PHP    *php.PHP
		MStore *mstore.MStore
	}
}

var (
	_ ingredient.WissKIFetcher = (*Prefixes)(nil)
)

// NoPrefix checks if this WissKI instance is excluded from generating prefixes.
// TODO: Move this to the database!
func (prefixes *Prefixes) NoPrefix() bool {
	liquid := ingredient.GetLiquid(prefixes)

	// FIXME: Ignoring error here!
	exists, _ := fsx.IsRegular(filepath.Join(liquid.FilesystemBase, "prefixes.skip"), false)
	return exists
}

//go:embed prefixes.php
var listURIPrefixesPHP string

// All returns the prefixes applying to this WissKI
//
// server is an optional server to fetch prefixes from.
// server may be nil.
func (prefixes *Prefixes) All(ctx context.Context, server *phpx.Server) ([]string, error) {
	uris, err := prefixes.getLivePrefixes(ctx, server)
	if err != nil {
		return nil, err
	}

	uris2, err := prefixes.filePrefixes()
	if err != nil {
		return nil, err
	}

	return append(uris, uris2...), nil
}

// getLivePrefixes get the list of prefixes found within the live system
func (prefixes *Prefixes) getLivePrefixes(ctx context.Context, server *phpx.Server) (pfs []string, err error) {
	danger := ingredient.GetStill(prefixes).Config.TS.DangerouslyUseAdapterPrefixes
	if !(danger.Set && danger.Value) {
		pfs, err = prefixes.getTSPrefixes(ctx, server)
	} else {
		// danger danger danger: Use the adapter prefixes
		pfs, err = prefixes.getAdapterPrefixes(ctx, server)
	}

	if err != nil {
		return nil, err
	}

	// sort the prefixes, and remove duplicates
	slices.Sort(pfs)
	pfs = collection.Deduplicate(pfs)

	// load the list of blocked prefixes
	blocks, err := prefixes.blocked()
	if err != nil {
		return nil, err
	}

	// filter out blocked prefixes
	return collection.Filter(pfs, func(uri string) bool { return !hasAnyPrefix(uri, blocks) }), nil
}

func (wisski *Prefixes) getAdapterPrefixes(ctx context.Context, server *phpx.Server) (pfs []string, err error) {
	err = wisski.dependencies.PHP.ExecScript(ctx, server, &pfs, listURIPrefixesPHP, "list_adapter_prefixes")
	if err != nil {
		return nil, err
	}
	return pfs, nil
}

func (wisski *Prefixes) getTSPrefixes(ctx context.Context, server *phpx.Server) (pfs []string, err error) {
	err = wisski.dependencies.PHP.ExecScript(ctx, server, &pfs, listURIPrefixesPHP, "list_triplestore_prefixes")
	if err != nil {
		return nil, err
	}
	return pfs, nil
}

func (prefixes *Prefixes) blocked() ([]string, error) {
	config := ingredient.GetStill(prefixes).Config

	// open the resolver block file
	// TODO: move this to the distillery
	file, err := os.Open(config.Paths.ResolverBlocks)
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
	return slices.ContainsFunc(
		prefixes,
		func(prefix string) bool {
			return strings.HasPrefix(candidate, prefix)
		},
	)
}

func (wisski *Prefixes) filePrefixes() (prefixes []string, err error) {
	path := filepath.Join(ingredient.GetLiquid(wisski).FilesystemBase, "prefixes")

	// check that the prefixes path exists
	{
		isFile, err := fsx.IsRegular(path, true)
		if err != nil {
			return nil, err
		}
		if !isFile {
			return nil, nil
		}
	}

	// open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// scan each line
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
func (wisski *Prefixes) AllCached(ctx context.Context) (results []string, err error) {
	return prefix.GetAll(ctx, wisski.dependencies.MStore)
}

// Update updates the cached prefixes of this instance
func (wisski *Prefixes) Update(ctx context.Context) error {
	prefixes, err := wisski.All(ctx, nil)
	if err != nil {
		return err
	}
	return prefix.SetAll(ctx, wisski.dependencies.MStore, prefixes...)
}

func (prefixes *Prefixes) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.NoPrefixes = prefixes.NoPrefix()
	if flags.Quick {
		// quick mode: grab only the cached prefixes
		info.Prefixes, _ = prefixes.AllCached(flags.Context)
	} else {
		// slow mode: grab the fresh prefixes from the server
		// TODO: Do we want to update them while we are at it?
		info.Prefixes, _ = prefixes.All(flags.Context, flags.Server)
	}
	return
}
