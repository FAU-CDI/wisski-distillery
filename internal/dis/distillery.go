// Package dis provides the main distillery
package dis

import (
	"io"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/api"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/next"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/panel"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/tokens"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/binder"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/docker"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/malt"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/purger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/resolver"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket/actions"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/cron"
	handleing "github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/home"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/legal"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/list"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/logo"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/news"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/solr"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2/sshkeys"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/web"
	"github.com/tkw1536/pkglib/lifetime"
)

// Distillery represents a WissKI Distillery
//
// It is the main structure used to interact with different components.
type Distillery struct {
	// core holds the core of the distillery
	component.Still

	// Where interactive progress is displayed
	Progress io.Writer

	// lifetime holds all components
	lifetime     lifetime.Lifetime[component.Component, component.Still]
	lifetimeInit sync.Once
}

//
// INIT & EXPORT
//

func (dis *Distillery) init() {
	dis.lifetimeInit.Do(func() {
		dis.lifetime.Init = func(c component.Component, s component.Still) {
			component.Init(c, s)
		}
		dis.lifetime.Register = dis.allComponents
	})
}

func export[C component.Component](dis *Distillery) C {
	dis.init()
	return lifetime.Export[C](&dis.lifetime, dis.Still)
}

func exportAll[C component.Component](dis *Distillery) []C {
	dis.init()
	return lifetime.ExportSlice[C](&dis.lifetime, dis.Still)
}

//
// PUBLIC COMPONENT GETTERS
//

func (dis *Distillery) Control() *server.Server {
	return export[*server.Server](dis)
}
func (dis *Distillery) SQL() *sql.SQL {
	return export[*sql.SQL](dis)
}
func (dis *Distillery) SSH() *ssh2.SSH2 {
	return export[*ssh2.SSH2](dis)
}
func (dis *Distillery) Auth() *auth.Auth {
	return export[*auth.Auth](dis)
}
func (dis *Distillery) Keys() *sshkeys.SSHKeys {
	return export[*sshkeys.SSHKeys](dis)
}
func (dis *Distillery) Cron() *cron.Cron {
	return export[*cron.Cron](dis)
}
func (dis *Distillery) Instances() *instances.Instances {
	return export[*instances.Instances](dis)
}
func (dis *Distillery) Exporter() *exporter.Exporter {
	return export[*exporter.Exporter](dis)
}
func (dis *Distillery) Provision() *provision.Provision {
	return export[*provision.Provision](dis)
}
func (dis *Distillery) Docker() *docker.Docker {
	return export[*docker.Docker](dis)
}

func (dis *Distillery) Installable() []component.Installable {
	return exportAll[component.Installable](dis)
}
func (dis *Distillery) Updatable() []component.Updatable {
	return exportAll[component.Updatable](dis)
}
func (dis *Distillery) Info() *admin.Admin {
	return export[*admin.Admin](dis)
}
func (dis *Distillery) Policy() *policy.Policy {
	return export[*policy.Policy](dis)
}
func (dis *Distillery) Templating() *templating.Templating {
	return export[*templating.Templating](dis)
}

func (dis *Distillery) Purger() *purger.Purger {
	return export[*purger.Purger](dis)
}

//
// All components
// THESE SHOULD NEVER BE CALLED DIRECTLY
//

func (dis *Distillery) allComponents(context *lifetime.Registry[component.Component, component.Still]) {
	lifetime.Place[*docker.Docker](context)
	lifetime.Place[*binder.Binder](context)
	lifetime.Place[*web.Web](context)

	lifetime.Register(context, func(ts *triplestore.Triplestore, _ component.Still) {
		ts.BaseURL = "http://" + dis.Upstream.TriplestoreAddr()
		ts.PollInterval = time.Second
	})

	lifetime.Register(context, func(sql *sql.SQL, _ component.Still) {
		sql.ServerURL = dis.Upstream.SQLAddr()
		sql.PollInterval = time.Second
	})
	lifetime.Place[*sql.LockTable](context)
	lifetime.Place[*sql.InstanceTable](context)

	lifetime.Register(context, func(s *solr.Solr, _ component.Still) {
		s.BaseURL = dis.Upstream.SolrAddr()
		s.PollInterval = time.Second
	})

	// auth
	lifetime.Place[*auth.Auth](context)
	lifetime.Place[*policy.Policy](context)
	lifetime.Place[*panel.UserPanel](context)
	lifetime.Place[*next.Next](context)
	lifetime.Place[*tokens.Tokens](context)

	//scopes
	lifetime.Place[*scopes.Never](context)
	lifetime.Place[*scopes.UserLoggedIn](context)
	lifetime.Place[*scopes.AdminLoggedIn](context)
	lifetime.Place[*scopes.ListInstancesScope](context)
	lifetime.Place[*scopes.ListNewsScope](context)
	lifetime.Place[*scopes.ResolverScope](context)

	// instances
	lifetime.Place[*instances.Instances](context)
	lifetime.Place[*meta.Meta](context)
	lifetime.Place[*malt.Malt](context)
	lifetime.Place[*provision.Provision](context)

	// Purger
	lifetime.Place[*purger.Purger](context)

	// Snapshots
	lifetime.Place[*exporter.Exporter](context)
	lifetime.Place[*logger.Logger](context)
	lifetime.Place[*exporter.Config](context)
	lifetime.Place[*exporter.Bookkeeping](context)
	lifetime.Place[*exporter.Filesystem](context)
	lifetime.Place[*exporter.Pathbuilders](context)

	// ssh server
	lifetime.Place[*ssh2.SSH2](context)
	lifetime.Place[*sshkeys.SSHKeys](context)

	// Control server
	lifetime.Place[*server.Server](context)
	lifetime.Place[*handleing.Handling](context)

	lifetime.Place[*home.Home](context)
	lifetime.Place[*list.ListInstances](context)
	lifetime.Register(context, func(resolver *resolver.Resolver, _ component.Still) {
		resolver.RefreshInterval = time.Minute
	})
	lifetime.Place[*admin.Admin](context) // TODO: Remove analytics
	lifetime.Place[*legal.Legal](context)
	lifetime.Place[*news.News](context)

	lifetime.Place[*assets.Static](context)
	lifetime.Place[*logo.Logo](context)
	lifetime.Place[*templating.Templating](context)

	// Websockets
	lifetime.Place[*socket.Sockets](context)
	lifetime.Place[*actions.Backup](context)
	lifetime.Place[*actions.Provision](context)
	lifetime.Place[*actions.Snapshot](context)
	lifetime.Place[*actions.Rebuild](context)
	lifetime.Place[*actions.Update](context)
	lifetime.Place[*actions.Cron](context)
	lifetime.Place[*actions.Start](context)
	lifetime.Place[*actions.Stop](context)
	lifetime.Place[*actions.Purge](context)

	// Cron
	lifetime.Place[*cron.Cron](context)

	// API
	lifetime.Place[*api.API](context)
	lifetime.Place[*list.API](context)
	lifetime.Place[*news.API](context)
	lifetime.Place[*resolver.API](context)
}
