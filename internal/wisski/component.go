package wisski

import (
	"path/filepath"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
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
	web lazy.Lazy[*web.Web]
	dis lazy.Lazy[*dis.Dis]
	ssh lazy.Lazy[*ssh.SSH]
	ts  lazy.Lazy[*triplestore.Triplestore]
	sql lazy.Lazy[*sql.SQL]

	// other components
	instances lazy.Lazy[*instances.Instances]
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
func makeComponent[C component.Component](dis *Distillery, field *lazy.Lazy[C], init func(C)) C {

	// get the typeof C and make sure that it is a pointer type!
	typC := reflect.TypeOf((*C)(nil)).Elem()
	if typC.Kind() != reflect.Pointer {
		panic("makeComponent: C must be backed by a pointer")
	}

	// return the field
	return field.Get(func() (c C) {
		c = reflect.New(typC.Elem()).Interface().(C)
		if init != nil {
			init(c)
		}

		base := c.Base()
		base.Config = dis.Config
		if base.Dir == "" {
			base.Dir = filepath.Join(dis.Config.DeployRoot, "core", c.Name())
		}

		return
	})
}

// Components returns all components that have a stack function
func (dis *Distillery) Components() []component.InstallableComponent {
	return []component.InstallableComponent{
		dis.Web(),
		dis.Dis(),
		dis.SSH(),
		dis.Triplestore(),
		dis.SQL(),
	}
}

func (dis *Distillery) Web() *web.Web {
	return makeComponent(dis, &dis.components.web, nil)
}

func (d *Distillery) Dis() *dis.Dis {
	return makeComponent(d, &d.components.dis, func(ddis *dis.Dis) {
		ddis.ResolverFile = core.PrefixConfig
		ddis.Instances = d.Instances()
	})
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
		instances.Dir = filepath.Join(dis.Config.DeployRoot, "instances")
		instances.SQL = dis.SQL()
		instances.TS = dis.Triplestore()
	})
}
