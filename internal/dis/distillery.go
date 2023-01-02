// Package dis provides the main distillery
package dis

import (
	"io"
	"sync"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/cron"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/home"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/info"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/malt"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/resolver"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/solr"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/web"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

// Distillery represents a WissKI Distillery
//
// It is the main structure used to interact with different components.
type Distillery struct {
	// core holds the core of the distillery
	component.Still

	// Where interactive progress is displayed
	Progress io.Writer

	// Upstream holds information to connect to the various running
	// distillery components.
	//
	// NOTE(twiesing): This is intended to eventually allow full remote management of the distillery.
	// But for now this will just hold upstream configuration.
	Upstream Upstream

	// pool holds all components
	pool     lazy.Pool[component.Component, component.Still]
	poolInit sync.Once
}

// Upstream contains the configuration for accessing remote configuration.
type Upstream struct {
	SQL         string
	Triplestore string
	Solr        string
}

//
// PUBLIC COMPONENT GETTERS
//

func (dis *Distillery) Control() *control.Control {
	return export[*control.Control](dis)
}
func (dis *Distillery) Resolver() *resolver.Resolver {
	return export[*resolver.Resolver](dis)
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

func (dis *Distillery) Cron() *cron.Cron {
	return export[*cron.Cron](dis)
}

func (dis *Distillery) Triplestore() *triplestore.Triplestore {
	return export[*triplestore.Triplestore](dis)
}
func (dis *Distillery) Instances() *instances.Instances {
	return export[*instances.Instances](dis)
}
func (dis *Distillery) Exporter() *exporter.Exporter {
	return export[*exporter.Exporter](dis)
}

func (dis *Distillery) Installable() []component.Installable {
	return exportAll[component.Installable](dis)
}
func (dis *Distillery) Updatable() []component.Updatable {
	return exportAll[component.Updatable](dis)
}
func (dis *Distillery) Provisionable() []component.Provisionable {
	return exportAll[component.Provisionable](dis)
}
func (dis *Distillery) Info() *info.Info {
	return export[*info.Info](dis)
}

//
// All components
// THESE SHOULD NEVER BE CALLED DIRECTLY
//

func (dis *Distillery) allComponents() []initFunc {
	return []initFunc{
		auto[*web.Web],

		manual(func(ts *triplestore.Triplestore) {
			ts.BaseURL = "http://" + dis.Upstream.Triplestore
			ts.PollInterval = time.Second
		}),
		manual(func(sql *sql.SQL) {
			sql.ServerURL = dis.Upstream.SQL
			sql.PollInterval = time.Second
		}),
		manual(func(s *solr.Solr) {
			s.BaseURL = dis.Upstream.Solr
			s.PollInterval = time.Second
		}),

		// auth
		auto[*auth.Auth],
		auto[*policy.Policy],

		// instances
		auto[*instances.Instances],
		auto[*meta.Meta],
		auto[*malt.Malt],

		// Snapshots
		auto[*exporter.Exporter],
		auto[*logger.Logger],
		auto[*exporter.Config],
		auto[*exporter.Bookkeeping],
		auto[*exporter.Filesystem],
		auto[*exporter.Pathbuilders],

		// ssh server
		auto[*ssh2.SSH2],

		// Control server
		auto[*control.Control],
		auto[*static.Static],
		auto[*home.Home],
		manual(func(resolver *resolver.Resolver) {
			resolver.RefreshInterval = time.Minute
		}),
		manual(func(info *info.Info) {
			info.Analytics = &dis.pool.Analytics
		}),

		// Cron
		auto[*cron.Cron],
		auto[*home.UpdateHome],
		auto[*home.UpdateInstanceList],
		auto[*home.UpdateRedirect],
	}
}
