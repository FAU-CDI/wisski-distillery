package dis

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/control"
	"github.com/FAU-CDI/wisski-distillery/internal/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/component/exporter/logger"

	"github.com/FAU-CDI/wisski-distillery/internal/component/home"
	"github.com/FAU-CDI/wisski-distillery/internal/component/info"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/component/resolver"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/component/ssh"
	"github.com/FAU-CDI/wisski-distillery/internal/component/static"
	"github.com/FAU-CDI/wisski-distillery/internal/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/component/web"
	"github.com/tkw1536/goprogram/lib/collection"
)

// register returns all components of the distillery
func (dis *Distillery) register(context component.ComponentPoolContext) []component.Component {
	return collection.MapSlice([]initFunc{
		auto[*web.Web],

		auto[*ssh.SSH],

		manual(func(ts *triplestore.Triplestore) {
			ts.BaseURL = "http://" + dis.Upstream.Triplestore
			ts.PollContext = dis.Context()
			ts.PollInterval = time.Second
		}),
		manual(func(sql *sql.SQL) {
			sql.ServerURL = dis.Upstream.SQL
			sql.PollContext = dis.Context()
			sql.PollInterval = time.Second
		}),

		auto[*instances.Instances],
		auto[*meta.Meta],

		// Snapshots
		auto[*exporter.Exporter],
		auto[*logger.Logger],
		auto[*exporter.Config],
		auto[*exporter.Bookkeeping],
		auto[*exporter.Filesystem],
		auto[*exporter.Pathbuilders],

		// Control server
		auto[*control.Control],
		auto[*static.Static],
		manual(func(home *home.Home) {
			home.RefreshInterval = time.Minute
		}),
		manual(func(resolver *resolver.Resolver) {
			resolver.RefreshInterval = time.Minute
		}),
		auto[*info.Info],
	}, func(f initFunc) component.Component {
		return f(dis, context)
	})
}

type initFunc = func(dis *Distillery, context component.ComponentPoolContext) component.Component

// manual initializes a component from the provided distillery.
func manual[C component.Component](init func(component C)) initFunc {
	return func(dis *Distillery, context component.ComponentPoolContext) component.Component {
		return component.MakeComponent(context, dis.Core, init)
	}
}

// use is like r, but does not provided additional initialization
func auto[C component.Component](dis *Distillery, context component.ComponentPoolContext) component.Component {
	return component.MakeComponent[C](context, dis.Core, nil)
}
