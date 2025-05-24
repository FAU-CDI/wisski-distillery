//spellchecker:words dockerx
package dockerx

//spellchecker:words context github docker types filters network client
import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// Client represents a docker client with additional functionality
type Client struct {
	*client.Client
}

// CreateNetwork creates a docker network with the given name unless it already exists.
// The new network will be of default type.
// exists indicates if the network already exists.
func (client *Client) NetworkCreate(ctx context.Context, name string) (id string, exists bool, e error) {
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
	create, err := client.Client.NetworkCreate(ctx, name, network.CreateOptions{Scope: "local"})
	if err != nil {
		return "", false, fmt.Errorf("failed to create docker network: %w", err)
	}
	return create.ID, false, nil
}
