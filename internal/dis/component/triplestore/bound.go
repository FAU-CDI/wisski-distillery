package triplestore

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
)

// For returns the bound triplestore for the given instance.
func (ts Triplestore) For(instance models.Instance) BoundTriplestore {
	// return either the global triplestore or the dedicated triplestore client.
	if !instance.DedicatedTriplestore {
		return &boundGlobal{
			client:   ts.globalClient(),
			instance: instance,
		}
	}

	return &boundDedicated{
		openStack: func() (*dockerx.Stack, error) {
			return dockerx.NewStack(ts.dependencies.Docker, instance.FilesystemBase)
		},
		serviceName: "dedicatedtriplestore",
		instance:    instance,
	}
}

type BoundTriplestore interface {
	// Returns the read and write URLs and credentials for WissKI instances accessing this triplestore.
	ReadURL() string
	WriteURL() string
	Credentials() (username string, password string)

	// Provisions or purges the repository belonging to this instance.
	Provision(ctx context.Context, domain string) error
	Purge(ctx context.Context, mustCreate bool) error

	// Snapshots or restores the repository belonging to this instance.
	SnapshotDB(ctx context.Context, dst io.Writer) error
	RestoreDB(ctx context.Context, reader io.Reader) error
}
