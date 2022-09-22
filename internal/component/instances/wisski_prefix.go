package instances

import (
	"bufio"
	"io"
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

// Prefixes returns the prefixes
func (wisski *WissKI) Prefixes() (prefixes []string, err error) {
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
	blocks, err := wisski.instances.BlockedPrefixes()
	if err != nil {
		return nil, err
	}

	// filter out blocked prefixes
	return slicesx.Filter(prefixes, func(uri string) bool { return !hasAnyPrefix(uri, blocks) }), nil
}

func (instances *Instances) BlockedPrefixes() ([]string, error) {
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

// PrefixConfig returns the prefix config belonging to this instance.
func (wisski *WissKI) PrefixConfig() (config string, err error) {
	// if the user requested to skip the prefix, then don't do anything with it!
	if wisski.NoPrefix() {
		return "", nil
	}

	var builder strings.Builder

	// domain
	builder.WriteString(wisski.URL().String() + ":")
	builder.WriteString("\n")

	// default prefixes
	prefixes, err := wisski.Prefixes()
	if err != nil {
		return "", err
	}

	// predefined prefixes
	for _, prefix := range prefixes {
		builder.WriteString(prefix)
		builder.WriteRune('\n')
	}

	// custom prefixes
	prefixPath := filepath.Join(wisski.FilesystemBase, "prefixes")
	if fsx.IsFile(wisski.instances.Environment, prefixPath) {
		prefix, err := wisski.instances.Core.Environment.Open(prefixPath)
		if err != nil {
			return "", err
		}
		defer prefix.Close()
		if _, err := io.Copy(&builder, prefix); err != nil {
			return "", err
		}
		builder.WriteString("\n")
	}

	// and done!
	return builder.String(), nil
}
