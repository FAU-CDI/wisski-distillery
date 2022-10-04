package dis

import (
	"reflect"
	"sync/atomic"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/control"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/component/resolver"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/component/ssh"
	"github.com/FAU-CDI/wisski-distillery/internal/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/component/web"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/slicesx"
	"github.com/tkw1536/goprogram/lib/reflectx"
)

// components holds the various components of the distillery
// It is inlined into the [Distillery] struct, and initialized using [makeComponent].
//
// The caller is responsible for syncronizing access across multiple goroutines.
type components struct {
	t    int32 // t is the previously used thread id!
	pool component.Pool
}

// register returns all components of the distillery
func register(dis *Distillery, thread int32) []component.Component {
	return []component.Component{
		ra[*web.Web](dis, thread),

		ra[*ssh.SSH](dis, thread),

		r(dis, thread, func(ts *triplestore.Triplestore) {
			ts.BaseURL = "http://" + dis.Upstream.Triplestore
			ts.PollContext = dis.Context()
			ts.PollInterval = time.Second
		}),
		r(dis, thread, func(sql *sql.SQL) {
			sql.ServerURL = dis.Upstream.SQL
			sql.PollContext = dis.Context()
			sql.PollInterval = time.Second
		}),

		ra[*instances.Instances](dis, thread),

		// Snapshots
		ra[*snapshots.Manager](dis, thread),
		ra[*snapshots.Config](dis, thread),
		ra[*snapshots.Bookkeeping](dis, thread),
		ra[*snapshots.Filesystem](dis, thread),
		ra[*snapshots.Pathbuilders](dis, thread),

		// Control server
		r(dis, thread, func(control *control.Control) {
			control.ResolverFile = core.PrefixConfig
		}),
		ra[*control.SelfHandler](dis, thread),
		r(dis, thread, func(resolver *resolver.Resolver) {
			resolver.ResolverFile = core.PrefixConfig
		}),
		ra[*control.Info](dis, thread),
	}
}

// r initializes a component from the provided distillery.
func r[C component.Component](dis *Distillery, thread int32, init func(component C)) C {
	return component.InitComponent(&dis.pool, thread, dis.Core, makeInitFunction(dis, init))
}

// ra is like r, but does not provided additional initialization
func ra[C component.Component](dis *Distillery, thread int32) C {
	return r[C](dis, thread, nil)
}

var componentType = reflectx.TypeOf[component.Component]()

// makeInitFunction generate an init function for a specific component.
// The function should be called at most once.
func makeInitFunction[C component.Component](dis *Distillery, rest func(instance C)) func(instance C, thread int32) {
	return func(instance C, thread int32) {
		meta := component.GetMeta[C]()

		// this function automatically initializes component.Component-like fields of the instance.
		// for this we first store two types of fields:

		singles := make(map[string]reflect.Type) // fields of type C where C is a component.Type
		multis := make(map[string]reflect.Type)  // fields of type []C where C is an interface that implements component.

		// fill the above variables, with a mapping of field name to struct
		count := meta.Elem.NumField()
		for i := 0; i < count; i++ {
			field := meta.Elem.Field(i)

			name := field.Name
			tp := field.Type

			switch {
			// field is a `Component``
			case tp.Implements(componentType):
				singles[name] = tp
			// field is a `[]Component``
			case tp.Kind() == reflect.Slice && tp.Elem().Kind() == reflect.Interface && tp.Elem().Implements(componentType):
				multis[name] = tp.Elem()
			}
		}

		// do the rest of the initialization
		if rest != nil {
			defer rest(instance)
		}

		if len(singles) == 0 && len(multis) == 0 {
			// no fields to assign, bail out immediatly
			return
		}

		registry := register(dis, thread)

		instanceV := reflect.ValueOf(instance).Elem()

		// assign the component fields
		for field, eType := range singles {
			c := slicesx.First(registry, func(c component.Component) bool {
				return reflect.TypeOf(c).AssignableTo(eType)
			})

			field := instanceV.FieldByName(field)
			field.Set(reflect.ValueOf(c))
		}

		// assign the multi subtypes
		registryR := reflect.ValueOf(registry)
		for field, eType := range multis {
			cs := filterSubtype(registryR, eType)
			field := instanceV.FieldByName(field)
			field.Set(cs)
		}
	}
}

// filterSubtype filters the slice of type []S into a slice of type []iface.
// S and I must be interface types.
func filterSubtype(sliceS reflect.Value, iface reflect.Type) reflect.Value {
	len := sliceS.Len()

	// convert each element
	result := reflect.MakeSlice(reflect.SliceOf(iface), 0, len)
	for i := 0; i < len; i++ {
		element := sliceS.Index(i)
		if element.Type().Implements(iface) {
			result = reflect.Append(result, element.Elem().Convert(iface))
		}
	}
	return result
}

//
// Export Components
//

// ea exports all components of the given subtype
func ea[C component.Component](dis *Distillery) []C {
	registry := register(dis, atomic.AddInt32(&dis.t, 1))

	results := make([]C, 0, len(registry))
	for _, comp := range registry {
		if cc, ok := comp.(C); ok {
			results = append(results, cc)
		}
	}
	return results
}

// e exports a single component of the given subtype
func e[C component.Component](dis *Distillery) C {
	for _, comp := range register(dis, atomic.AddInt32(&dis.t, 1)) {
		if cc, ok := comp.(C); ok {
			return cc
		}
	}
	panic("e: component is missing")
}
