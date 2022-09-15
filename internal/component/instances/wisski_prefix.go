package instances

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/tkw1536/goprogram/stream"
)

// NoPrefix checks if this WissKI instance is excluded from generating prefixes.
// TODO: Move this to the database!
func (wisski *WissKI) NoPrefix() bool {
	return fsx.IsFile(filepath.Join(wisski.FilesystemBase, "prefixes.skip"))
}

var errPrefixExecFailed = errors.New("PrefixConfig: Failed to call list_uri_prefixes")

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
	wu := stream.NewIOStream(&builder, nil, nil, 0)
	code, err := wisski.Barrel().Exec(wu, "barrel", "/bin/bash", "/user_shell.sh", "-c", "drush php:script /wisskiutils/list_uri_prefixes.php")
	if err != nil || code != 0 {
		return "", errPrefixExecFailed
	}

	// custom prefixes
	prefixPath := filepath.Join(wisski.FilesystemBase, "prefixes")
	if fsx.IsFile(prefixPath) {
		prefix, err := os.Open(prefixPath)
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
