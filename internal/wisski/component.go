package wisski

import (
	"path/filepath"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/control"
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
	web     lazy.Lazy[*web.Web]
	control lazy.Lazy[*control.Control]
	ssh     lazy.Lazy[*ssh.SSH]
	ts      lazy.Lazy[*triplestore.Triplestore]
	sql     lazy.Lazy[*sql.SQL]

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
		base.Core = dis.Core
		if base.Dir == "" {
			base.Dir = filepath.Join(dis.Config.DeployRoot, "core", c.Name())
		}

		return
	})
}

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

// Backupable returns all the components that can be backuped up.
func (dis *Distillery) Backupable() []component.Backupable {
	return getComponents[component.Backupable](dis)
}

// Installables returns all components that can be installed
func (dis *Distillery) Installables() []component.Installable {
	return getComponents[component.Installable](dis)
}

// Installables returns all components that can be installed
func (dis *Distillery) Updateable() []component.Updatable {
	return getComponents[component.Updatable](dis)
}

// Provisionable returns all components which can be provisioned
func (dis *Distillery) Provisionable() []component.Provisionable {
	return getComponents[component.Provisionable](dis)
}

func getComponents[C component.Component](dis *Distillery) (result []C) {
	all := dis.Components()

	result = make([]C, 0, len(all))
	for _, c := range all {
		sc, ok := c.(C)
		if !ok {
			continue
		}
		result = append(result, sc)
	}

	return
}

func (dis *Distillery) Web() *web.Web {
	return makeComponent(dis, &dis.components.web, nil)
}

func (d *Distillery) Control() *control.Control {
	return makeComponent(d, &d.components.control, func(ddis *control.Control) {
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
