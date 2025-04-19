//spellchecker:words docker
package docker

//spellchecker:words context github wisski distillery internal component docker types filters client
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Docker struct {
	component.Base
}

// DockerClient is a client to the docker api
type DockerClient = *client.Client

func (docker *Docker) APIClient() (DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// Ping pings the docker daemon to check if it is properly working
func (docker *Docker) Ping(ctx context.Context) (types.Ping, error) {
	client, err := docker.APIClient()
	if err != nil {
		return types.Ping{}, err
	}
	defer client.Close()

	ping, err := client.Ping(ctx)
	if err != nil {
		return types.Ping{}, err
	}

	return ping, err
}

// CreateNetwork creates a docker network with the given name unless it already exists.
// The new network will be of default type.
// exists indicates if the network already exists.
func (docker *Docker) CreateNetwork(ctx context.Context, name string) (id string, exists bool, err error) {
	// create a new docker client
	client, err := docker.APIClient()
	if err != nil {
		return "", false, err
	}
	defer client.Close()

	// check if the network exists, by listing the network name
	list, err := client.NetworkList(ctx, network.ListOptions{Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: name})})
	if err != nil {
		return "", false, err
	}

	// network already exists => nothing to do
	if len(list) == 1 {
		return list[0].ID, true, nil
	}

	// do the actual create!
	create, err := client.NetworkCreate(ctx, name, network.CreateOptions{Scope: "local"})
	if err != nil {
		return "", false, err
	}
	return create.ID, false, nil
}
