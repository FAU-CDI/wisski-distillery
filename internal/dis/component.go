package dis

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/control"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/component/ssh"
	"github.com/FAU-CDI/wisski-distillery/internal/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/component/web"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

// components holds the various components of the distillery
// It is inlined into the [Distillery] struct, and initialized using [makeComponent].
//
// The caller is responsible for syncronizing access across multiple goroutines.
type components struct {

	// installable components
	web     lazy.Lazy[*web.Web]
	control lazy.Lazy[*control.Control]
	ssh     lazy.Lazy[*ssh.SSH]
	ts      lazy.Lazy[*triplestore.Triplestore]
	sql     lazy.Lazy[*sql.SQL]

	// other components
	instances lazy.Lazy[*instances.Instances]
	snapshots lazy.Lazy[*snapshots.Manager]
}

//
// Individual Components
//

func (dis *Distillery) Web() *web.Web {
	return component.Initialize(dis.Core, &dis.components.web, nil)
}

func (d *Distillery) Control() *control.Control {
	return component.Initialize(d.Core, &d.components.control, func(control *control.Control) {
		control.ResolverFile = core.PrefixConfig
		control.Instances = d.Instances()
	})
}

func (dis *Distillery) SSH() *ssh.SSH {
	return component.Initialize(dis.Core, &dis.components.ssh, nil)
}

func (dis *Distillery) SQL() *sql.SQL {
	return component.Initialize(dis.Core, &dis.components.sql, func(sql *sql.SQL) {
		sql.ServerURL = dis.Upstream.SQL
		sql.PollContext = dis.Context()
		sql.PollInterval = time.Second
	})
}

func (dis *Distillery) Triplestore() *triplestore.Triplestore {
	return component.Initialize(dis.Core, &dis.components.ts, func(ts *triplestore.Triplestore) {
		ts.BaseURL = "http://" + dis.Upstream.Triplestore
		ts.PollContext = dis.Context()
		ts.PollInterval = time.Second
	})
}

func (dis *Distillery) Instances() *instances.Instances {
	return component.Initialize(dis.Core, &dis.components.instances, func(instances *instances.Instances) {
		instances.SQL = dis.SQL()
		instances.TS = dis.Triplestore()
	})
}

func (dis *Distillery) SnapshotManager() *snapshots.Manager {
	return component.Initialize(dis.Core, &dis.components.snapshots, func(snapshots *snapshots.Manager) {
		snapshots.SQL = dis.SQL()
		snapshots.TS = dis.Triplestore()
	})
}

//
// ALL COMPONENTS
//

func (dis *Distillery) Components() []component.Component {
	return []component.Component{
		dis.Web(),
		dis.Control(),
		dis.SSH(),
		dis.Triplestore(),
		dis.SQL(),
		dis.Instances(),
	}
}

//
// COMPONENT SUBTYPE GETTERS
//

// Backupable returns all the components that can be backuped up.
func (dis *Distillery) Backupable() []component.Backupable {
	return getComponentSubtype[component.Backupable](dis)
}

// Installables returns all components that can be installed
func (dis *Distillery) Installables() []component.Installable {
	return getComponentSubtype[component.Installable](dis)
}

// Installables returns all components that can be installed
func (dis *Distillery) Updateable() []component.Updatable {
	return getComponentSubtype[component.Updatable](dis)
}

// Provisionable returns all components which can be provisioned
func (dis *Distillery) Provisionable() []component.Provisionable {
	return getComponentSubtype[component.Provisionable](dis)
}

// getComponentSubtype gets all components of type T
func getComponentSubtype[T component.Component](dis *Distillery) (components []T) {
	all := dis.Components()

	components = make([]T, 0, len(all))
	for _, c := range all {
		sc, ok := c.(T)
		if !ok {
			continue
		}
		components = append(components, sc)
	}

	return
}
