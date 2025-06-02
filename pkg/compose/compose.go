//spellchecker:words compose
package compose

//spellchecker:words path filepath github compose spec loader types
import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
)

// ComposeProject represents a compose project.
type Project = *types.Project

// Open loads a docker compose project from the given path.
// The filename of the compose file and environment file names are assumed to be default.
func Open(ctx context.Context, path string) (project Project, err error) {
	wd, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve working directory for compose: %w", err)
	}

	opts, err := cli.NewProjectOptions(
		nil,

		cli.WithWorkingDirectory(wd),
		cli.WithEnvFiles(),
		cli.WithDotEnv,
		cli.WithConfigFileEnv,
		cli.WithDefaultConfigPath,

		cli.WithInterpolation(true),
		cli.WithResolvedPaths(true),
		cli.WithNormalization(true),
		cli.WithConsistency(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project options: %w", err)
	}

	proj, err := cli.ProjectFromOptions(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create compose project: %w", err)
	}

	return proj, nil
}
