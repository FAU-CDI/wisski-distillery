//spellchecker:words docker
package docker

//spellchecker:words github wisski distillery internal component dockerx docker client
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"github.com/docker/docker/client"
)

// Docker implements [dockerx.Factory]
type Docker struct {
	component.Base
}

func (docker *Docker) NewClient() (*dockerx.Client, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate client: %w", err)
	}
	return &dockerx.Client{Client: client}, nil
}
