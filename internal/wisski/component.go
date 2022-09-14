package wisski

import (
	"path/filepath"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/resolver"
	"github.com/FAU-CDI/wisski-distillery/internal/component/self"
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

	// installable components
	web      *web.Web
	self     *self.Self
	resolver *resolver.Resolver
	dis      *dis.Dis
	ssh      *ssh.SSH
	ts       *triplestore.Triplestore
	sql      *sql.SQL

	// other components
	instances *instances.Instances
}

// makeComponent makes or returns a component inside the [component] struct of the distillery
//
// C is the type of component to initialize. It must be backed by a pointer, or makeComponent will panic.
//
// dis is the distillery to initialize components for
// field is a pointer to the appropriate struct field within the distillery components
// init is called with a new non-nil component to initialize it. It may be nil, to indicate no initialization is required.
//
// makeComponent returns the new or existing component instance
func makeComponent[C component.Component](dis *Distillery, field *C, init func(C)) C {

	// get the typeof C and make sure that it is a pointer type!
	typC := reflect.TypeOf((*C)(nil)).Elem()
	if typC.Kind() != reflect.Pointer {
		panic("makeComponent: C must be backed by a pointer")
	}

	// if the component is non-nil, then it has already been initialized
	if !reflect.ValueOf(*field).IsNil() {
		return *field
	}

	// create a new element, and call the initializer (if requested)
	*field = reflect.New(typC.Elem()).Interface().(C)
	if init != nil {
		init(*field)
	}

	// apply the base configuration
	base := (*field).Base()
	base.Config = dis.Config
	base.Dir = filepath.Join(dis.Config.DeployRoot, "core", (*field).Name())

	// and eventually return it
	return *field
}

// Components returns all components that have a stack function
func (dis *Distillery) Components() []component.InstallableComponent {
	return []component.InstallableComponent{
		dis.Web(),
		dis.Self(),
		dis.Resolver(),
		dis.Dis(),
		dis.SSH(),
		dis.Triplestore(),
		dis.SQL(),
	}
}

func (dis *Distillery) Web() *web.Web {
	return makeComponent(dis, &dis.components.web, nil)
}

func (dis *Distillery) Self() *self.Self {
	return makeComponent(dis, &dis.components.self, nil)
}

func (dis *Distillery) Resolver() *resolver.Resolver {
	return makeComponent(dis, &dis.components.resolver, func(resolver *resolver.Resolver) {
		resolver.ConfigName = core.PrefixConfig
	})
}

func (d *Distillery) Dis() *dis.Dis {
	return makeComponent(d, &d.components.dis, nil)
}

func (dis *Distillery) SSH() *ssh.SSH {
	return makeComponent(dis, &dis.components.ssh, nil)
}

func (dis *Distillery) SQL() *sql.SQL {
	return makeComponent(dis, &dis.components.sql, func(sql *sql.SQL) {
		sql.ServerURL = dis.Upstream.SQL
		sql.PollContext = dis.Context()
		sql.PollInterval = time.Second
	})
}

func (dis *Distillery) Triplestore() *triplestore.Triplestore {
	return makeComponent(dis, &dis.components.ts, func(ts *triplestore.Triplestore) {
		ts.BaseURL = "http://" + dis.Upstream.Triplestore
		ts.PollContext = dis.Context()
		ts.PollInterval = time.Second
	})
}

func (dis *Distillery) Instances() *instances.Instances {
	return makeComponent(dis, &dis.components.instances, func(instances *instances.Instances) {
		instances.SQL = dis.SQL()
		instances.TS = dis.Triplestore()
	})
}
