package env

import (
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/component"
	"github.com/FAU-CDI/wisski-distillery/component/dis"
	"github.com/FAU-CDI/wisski-distillery/component/resolver"
	"github.com/FAU-CDI/wisski-distillery/component/self"
	"github.com/FAU-CDI/wisski-distillery/component/sql"
	"github.com/FAU-CDI/wisski-distillery/component/ssh"
	"github.com/FAU-CDI/wisski-distillery/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/component/web"
)

// TODO: Remove me when migration is complete
type Component = component.Component

// TODO: Move everything into specific subpackages

// Stacks returns the Stacks of this distillery
func (dis *Distillery) Components() []component.Component {
	// TODO: Do we want to cache these components?
	return []Component{
		dis.Web(),
		dis.Self(),
		dis.Resolver(),
		dis.Dis(),
		dis.SSH(),
		dis.Triplestore(),
		dis.SQL(),
	}
}

// Web returns the web component belonging to this distillery
func (dis *Distillery) Web() (web web.Web) {
	dis.makeComponent(web, &web.ComponentBase)
	return
}

// Self returns the self component belonging to this distillery
func (dis *Distillery) Self() (self self.Self) {
	dis.makeComponent(self, &self.ComponentBase)
	return
}

// Resolver returns the resolver component belonging to this distillery
func (dis *Distillery) Resolver() (resolver resolver.Resolver) {
	resolver.ConfigName = "prefix.cfg" // TODO: Move into core?
	resolver.Executable = dis.CurrentExecutable()

	dis.makeComponent(resolver, &resolver.ComponentBase)
	return
}

// Dis returns the dis component belonging to this distillery
func (dis *Distillery) Dis() (ddis dis.Dis) {
	ddis.Executable = dis.CurrentExecutable()

	dis.makeComponent(ddis, &ddis.ComponentBase)
	return
}

// SSH returns the SSH component belonging to this distillery
func (dis *Distillery) SSH() (ssh ssh.SSH) {
	dis.makeComponent(ssh, &ssh.ComponentBase)
	return
}

// SQL returns the SQL component belonging to this distillery
func (dis *Distillery) SQL() (sql sql.SQL) {
	sql.ServerURL = dis.Upstream.SQL
	sql.PollContext = dis.Context()
	sql.PollInterval = time.Second

	dis.makeComponent(sql, &sql.ComponentBase)
	return
}

// Triplestore returns the TriplestoreComponent belonging to this distillery
func (dis *Distillery) Triplestore() (ts triplestore.Triplestore) {
	ts.BaseURL = "http://" + dis.Upstream.Triplestore
	ts.PollContext = dis.Context()
	ts.PollInterval = time.Second

	dis.makeComponent(ts, &ts.ComponentBase)
	return
}

// makeComponent updates the baseComponent belonging to component
func (dis *Distillery) makeComponent(component component.Component, base *component.ComponentBase) {
	base.Config = dis.Config
	base.Dir = filepath.Join(dis.Config.DeployRoot, "core", component.Name())
}
