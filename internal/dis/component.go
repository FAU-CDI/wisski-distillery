package dis

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/control"
	"github.com/FAU-CDI/wisski-distillery/internal/component/home"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/resolver"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/component/ssh"
	"github.com/FAU-CDI/wisski-distillery/internal/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/component/web"
)

// register returns all components of the distillery
func (dis *Distillery) register(context *component.PoolContext) []component.Component {
	return []component.Component{
		ra[*web.Web](dis, context),

		ra[*ssh.SSH](dis, context),

		r(dis, context, func(ts *triplestore.Triplestore) {
			ts.BaseURL = "http://" + dis.Upstream.Triplestore
			ts.PollContext = dis.Context()
			ts.PollInterval = time.Second
		}),
		r(dis, context, func(sql *sql.SQL) {
			sql.ServerURL = dis.Upstream.SQL
			sql.PollContext = dis.Context()
			sql.PollInterval = time.Second
		}),

		ra[*instances.Instances](dis, context),

		// Snapshots
		ra[*snapshots.Manager](dis, context),
		ra[*snapshots.Config](dis, context),
		ra[*snapshots.Bookkeeping](dis, context),
		ra[*snapshots.Filesystem](dis, context),
		ra[*snapshots.Pathbuilders](dis, context),

		// Control server
		ra[*control.Control](dis, context),
		r(dis, context, func(home *home.Home) {
			home.RefreshInterval = time.Minute
		}),
		r(dis, context, func(resolver *resolver.Resolver) {
			resolver.RefreshInterval = time.Minute
		}),
		ra[*control.Info](dis, context),
	}
}

// r initializes a component from the provided distillery.
func r[C component.Component](dis *Distillery, context *component.PoolContext, init func(component C)) C {
	return component.Make(context, dis.Core, init)
}

// ra is like r, but does not provided additional initialization
func ra[C component.Component](dis *Distillery, context *component.PoolContext) C {
	return r[C](dis, context, nil)
}
