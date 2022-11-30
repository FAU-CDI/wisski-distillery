package lazy

import (
	"reflect"

	"github.com/tkw1536/goprogram/lib/collection"
	"golang.org/x/exp/slices"
)

type PoolAnalytics struct {
	Components map[string]*PoolAnalyticsComponent
	Groups     map[string]*PoolAnalyticsGroup
}

type PoolAnalyticsComponent struct {
	Type   string   // Type name
	Groups []string // groups this is contained in

	CFields map[string]string // fields with type C for which C implements component
	IFields map[string]string // fields []I where I is an interface that implements component

	Methods map[string]string // Method signatures of type
}
type PoolAnalyticsGroup struct {
	Type       string   // Type name
	Components []string // Components of this Type

	Methods map[string]string // Method signatures of this interface
}

// anal writes analytics about this context to anal
func (context *PoolContext[Component]) anal(anal *PoolAnalytics, groups []reflect.Type) {
	anal.Components = make(map[string]*PoolAnalyticsComponent, len(context.metaCache))
	anal.Groups = make(map[string]*PoolAnalyticsGroup)

	// collect all the pointers, and setup the anal.Components map!
	tpPointers := make([]reflect.Type, 0, len(context.metaCache))
	for _, meta := range context.metaCache {
		tp := reflect.PointerTo(meta.Elem)
		tpPointers = append(tpPointers, tp)

		mcount := tp.NumMethod()

		anal.Components[meta.Name] = &PoolAnalyticsComponent{
			Groups:  make([]string, 0),
			Methods: make(map[string]string, mcount),
		}
		for i := 0; i < mcount; i++ {
			method := tp.Method(i)
			anal.Components[meta.Name].Methods[method.Name] = method.Type.String()
		}
	}

	// collect interfaces to analyze
	ifaces := make([]reflect.Type, len(groups))
	copy(ifaces, groups)

	// take all of the components out of the cache
	for _, meta := range context.metaCache {
		anal.Components[meta.Name].Type = meta.Name
		anal.Components[meta.Name].CFields = collection.MapValues(meta.CFields, func(key string, tp reflect.Type) string {
			return nameOf(tp.Elem())
		})

		anal.Components[meta.Name].IFields = collection.MapValues(meta.IFields, func(key string, iface reflect.Type) string {
			ifaces = append(ifaces, iface)
			return nameOf(iface)
		})
	}

	// and analyze all ifaces
	for _, iface := range ifaces {
		name := nameOf(iface)
		if _, ok := anal.Groups[name]; ok {
			continue
		}

		types := collection.FilterClone(tpPointers, func(tp reflect.Type) bool {
			return tp.AssignableTo(iface)
		})

		anal.Groups[name] = &PoolAnalyticsGroup{
			Type: name,
			Components: collection.MapSlice(types, func(tp reflect.Type) string {
				cname := nameOf(tp.Elem())
				anal.Components[cname].Groups = append(anal.Components[cname].Groups, name)
				return cname
			}),
		}

		mcount := iface.NumMethod()
		anal.Groups[name].Methods = make(map[string]string, mcount)
		for i := 0; i < mcount; i++ {
			method := iface.Method(i)
			anal.Groups[name].Methods[method.Name] = method.Type.String()
		}
	}

	for _, comp := range anal.Components {
		slices.Sort(comp.Groups)
	}
}
