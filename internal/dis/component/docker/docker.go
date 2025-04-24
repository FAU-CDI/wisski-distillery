//spellchecker:words docker
package docker

//spellchecker:words context github wisski distillery internal component docker types filters client
import (
	"context"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/errwrap"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Docker struct {
	component.Base
}

// DockerClient is a client to the docker api.
// TODO: Make this private
type DockerClient = *client.Client

// APIClient creates a new docker api client.
// The caller must close the client.
func (docker *Docker) APIClient() (DockerClient, error) {
	// TODO: make this function private?
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate client: %w", err)
	}
	return cli, nil
}

// Ping pings the docker daemon to check if it is properly working.
func (docker *Docker) Ping(ctx context.Context) (p types.Ping, e error) {
	client, err := docker.APIClient()
	if err != nil {
		return types.Ping{}, fmt.Errorf("failed to create docker client: %w", err)
	}
	defer errwrap.Close(client, "docker client", &e)

	ping, err := client.Ping(ctx)
	if err != nil {
		return types.Ping{}, fmt.Errorf("failed to ping docker daemon: %w", err)
	}

	return ping, nil
}

// CreateNetwork creates a docker network with the given name unless it already exists.
// The new network will be of default type.
// exists indicates if the network already exists.
func (docker *Docker) CreateNetwork(ctx context.Context, name string) (id string, exists bool, e error) {
	// create a new docker client
	client, err := docker.APIClient()
	if err != nil {
		return "", false, fmt.Errorf("failed to create docker client: %w", err)
	}
	defer errwrap.Close(client, "docker client", &e)

	// check if the network exists, by listing the network name
	list, err := client.NetworkList(ctx, network.ListOptions{Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: name})})
	if err != nil {
		return "", false, fmt.Errorf("failed to list docker networks: %w", err)
	}

	// network already exists => nothing to do
	if len(list) == 1 {
		return list[0].ID, true, nil
	}

	// do the actual create!
	create, err := client.NetworkCreate(ctx, name, network.CreateOptions{Scope: "local"})
	if err != nil {
		return "", false, fmt.Errorf("failed to create docker network: %w", err)
	}
	return create.ID, false, nil
}
