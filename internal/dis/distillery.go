// Package dis provides the main distillery
package dis

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/control"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/resolver"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/component/ssh"
	"github.com/FAU-CDI/wisski-distillery/internal/component/triplestore"
)

// Distillery represents a WissKI Distillery
//
// It is the main structure used to interact with different components.
type Distillery struct {
	// core holds the core of the distillery
	component.Core

	// Upstream holds information to connect to the various running
	// distillery components.
	//
	// NOTE(twiesing): This is intended to eventually allow full remote management of the distillery.
	// But for now this will just hold upstream configuration.
	Upstream Upstream

	// components hold references to the various components of the distillery.
	components
}

// Upstream contains the configuration for accessing remote configuration.
type Upstream struct {
	SQL         string
	Triplestore string
}

// Context returns a new Context belonging to this distillery
func (dis *Distillery) Context() context.Context {
	return context.Background()
}

//
// PUBLIC COMPONENT GETTERS
//

func (dis *Distillery) Control() *control.Control {
	return dis.cControl(dis.thread())
}
func (dis *Distillery) Resolver() *resolver.Resolver {
	return dis.cResolver(dis.thread())
}
func (dis *Distillery) SSH() *ssh.SSH {
	return dis.cSSH(dis.thread())
}
func (dis *Distillery) SQL() *sql.SQL {
	return dis.cSQL(dis.thread())
}
func (dis *Distillery) Triplestore() *triplestore.Triplestore {
	return dis.cTriplestore(dis.thread())
}
func (dis *Distillery) Instances() *instances.Instances {
	return dis.cInstances(dis.thread())
}
func (dis *Distillery) SnapshotManager() *snapshots.Manager {
	return dis.cSnapshotManager(dis.thread())
}

func (dis *Distillery) Installable() []component.Installable { return dis.cInstallables(dis.thread()) }
func (dis *Distillery) Updatable() []component.Updatable     { return dis.cUpdateable(dis.thread()) }

func (dis *Distillery) Provisionable() []component.Provisionable {
	return dis.cProvisionable(dis.thread())
}
