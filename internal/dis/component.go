package dis

import (
	"sync/atomic"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/control"
	"github.com/FAU-CDI/wisski-distillery/internal/component/extras"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/component/ssh"
	"github.com/FAU-CDI/wisski-distillery/internal/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/component/web"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
)

// components holds the various components of the distillery
// It is inlined into the [Distillery] struct, and initialized using [makeComponent].
//
// The caller is responsible for syncronizing access across multiple goroutines.
type components struct {
	t    int32 // t is the previously used thread id!
	pool component.Pool
}

//
// Individual Components
//

func (c *components) thread() int32 {
	return atomic.AddInt32(&c.t, 1)
}

func (dis *Distillery) cWeb(thread int32) *web.Web {
	return component.PutComponent[*web.Web](&dis.pool, thread, dis.Core, nil)
}

func (dis *Distillery) cControl(thread int32) *control.Control {
	return component.PutComponent(&dis.pool, thread, dis.Core, func(control *control.Control, thread int32) {
		control.ResolverFile = core.PrefixConfig
		control.Instances = dis.cInstances(thread)
	})
}

func (dis *Distillery) cSSH(thread int32) *ssh.SSH {
	return component.PutComponent[*ssh.SSH](&dis.pool, thread, dis.Core, nil)
}

func (dis *Distillery) cSQL(thread int32) *sql.SQL {
	return component.PutComponent(&dis.pool, thread, dis.Core, func(sql *sql.SQL, thread int32) {
		sql.ServerURL = dis.Upstream.SQL
		sql.PollContext = dis.Context()
		sql.PollInterval = time.Second
	})
}

func (dis *Distillery) cTriplestore(thread int32) *triplestore.Triplestore {
	return component.PutComponent(&dis.pool, thread, dis.Core, func(ts *triplestore.Triplestore, thread int32) {
		ts.BaseURL = "http://" + dis.Upstream.Triplestore
		ts.PollContext = dis.Context()
		ts.PollInterval = time.Second
	})
}

func (dis *Distillery) cInstances(thread int32) *instances.Instances {
	return component.PutComponent(&dis.pool, thread, dis.Core, func(instances *instances.Instances, thread int32) {
		instances.SQL = dis.cSQL(thread)
		instances.TS = dis.cTriplestore(thread)
	})
}

func (dis *Distillery) cSnapshotManager(thread int32) *snapshots.Manager {
	return component.PutComponent(&dis.pool, thread, dis.Core, func(snapshots *snapshots.Manager, thread int32) {
		snapshots.SQL = dis.cSQL(thread)
		snapshots.Instances = dis.cInstances(thread)
		snapshots.Snapshotable = dis.cSnapshotable(thread)
		snapshots.Backupable = dis.cBackupable(thread)
	})
}

//
// EXTRAS COMPONENTS
//

func (dis *Distillery) cExtrasConfig(thread int32) *extras.Config {
	return component.PutComponent[*extras.Config](&dis.pool, thread, dis.Core, nil)
}

func (dis *Distillery) cExtrasBookkeeping(thread int32) *extras.Bookkeeping {
	return component.PutComponent[*extras.Bookkeeping](&dis.pool, thread, dis.Core, nil)
}

func (dis *Distillery) cExtrasFilesystem(thread int32) *extras.Filesystem {
	return component.PutComponent[*extras.Filesystem](&dis.pool, thread, dis.Core, nil)
}

func (dis *Distillery) cExtrasPathbuilders(thread int32) *extras.Pathbuilders {
	return component.PutComponent(&dis.pool, thread, dis.Core, func(pbs *extras.Pathbuilders, thread int32) {
		pbs.Instances = dis.cInstances(thread)
	})
}

//
// ALL COMPONENTS
//

func (dis *Distillery) cComponents(thread int32) []component.Component {
	return []component.Component{
		dis.cWeb(thread),
		dis.cControl(thread),
		dis.cSSH(thread),
		dis.cTriplestore(thread),
		dis.cSQL(thread),
		dis.cInstances(thread),
		dis.cSnapshotManager(thread),

		dis.cExtrasConfig(thread),
		dis.cExtrasBookkeeping(thread),
		dis.cExtrasFilesystem(thread),
		dis.cExtrasPathbuilders(thread),
	}
}

//
// COMPONENT SUBTYPE GETTERS
//

func (dis *Distillery) cInstallables(thread int32) []component.Installable {
	return getComponentSubtype[component.Installable](dis, thread)
}

func (dis *Distillery) cUpdateable(thread int32) []component.Updatable {
	return getComponentSubtype[component.Updatable](dis, thread)
}

func (dis *Distillery) cBackupable(thread int32) []component.Backupable {
	return getComponentSubtype[component.Backupable](dis, thread)
}

func (dis *Distillery) cProvisionable(thread int32) []component.Provisionable {
	return getComponentSubtype[component.Provisionable](dis, thread)
}

func (dis *Distillery) cSnapshotable(thread int32) []component.Snapshotable {
	return getComponentSubtype[component.Snapshotable](dis, thread)
}

func getComponentSubtype[T component.Component](dis *Distillery, thread int32) (components []T) {
	all := dis.cComponents(thread)

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
