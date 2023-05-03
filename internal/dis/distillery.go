// Package dis provides the main distillery
package dis

import (
	"io"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/next"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/panel"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
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
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/cron"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/home"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/legal"
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

func (dis *Distillery) allComponents() []initFunc {
	return []initFunc{
		auto[*docker.Docker],
		auto[*binder.Binder],
		auto[*web.Web],

		manual(func(ts *triplestore.Triplestore) {
			ts.BaseURL = "http://" + dis.Upstream.TriplestoreAddr()
			ts.PollInterval = time.Second
		}),

		manual(func(sql *sql.SQL) {
			sql.ServerURL = dis.Upstream.SQLAddr()
			sql.PollInterval = time.Second
		}),
		auto[*sql.LockTable],
		auto[*sql.InstanceTable],

		manual(func(s *solr.Solr) {
			s.BaseURL = dis.Upstream.SolrAddr()
			s.PollInterval = time.Second
		}),

		// auth
		auto[*auth.Auth],
		auto[*policy.Policy],
		auto[*panel.UserPanel],
		auto[*next.Next],

		//scopes
		auto[*scopes.UserLoggedIn],
		auto[*scopes.AdminLoggedIn],

		// instances
		auto[*instances.Instances],
		auto[*meta.Meta],
		auto[*malt.Malt],
		auto[*provision.Provision],

		// Purger
		auto[*purger.Purger],

		// Snapshots
		auto[*exporter.Exporter],
		auto[*logger.Logger],
		auto[*exporter.Config],
		auto[*exporter.Bookkeeping],
		auto[*exporter.Filesystem],
		auto[*exporter.Pathbuilders],

		// ssh server
		auto[*ssh2.SSH2],
		auto[*sshkeys.SSHKeys],

		// Control server
		auto[*server.Server],

		auto[*home.Home],
		manual(func(resolver *resolver.Resolver) {
			resolver.RefreshInterval = time.Minute
		}),
		manual(func(admin *admin.Admin) {
			admin.Analytics = &dis.lifetime.Analytics
		}),
		auto[*socket.Sockets],
		auto[*legal.Legal],
		auto[*news.News],

		auto[*assets.Static],
		auto[*logo.Logo],
		auto[*templating.Templating],

		// Cron
		auto[*cron.Cron],
		auto[*home.UpdateHome],
		auto[*home.UpdateInstanceList],

		// API
		auto[*home.API],
		auto[*news.API],
	}
}
