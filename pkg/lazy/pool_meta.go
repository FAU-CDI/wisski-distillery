package lazy

import (
	"reflect"

	"github.com/tkw1536/goprogram/lib/collection"
	"github.com/tkw1536/goprogram/lib/reflectx"
)

// getMeta gets the component belonging to a component type
func getMeta[Component any, ConcreteComponent any](cache map[reflect.Type]meta[Component]) meta[Component] {
	tp := reflectx.TypeOf[ConcreteComponent]()

	// we already have a m => return it
	if m, ok := cache[tp]; ok {
		return m
	}

	// create a new m
	var m meta[Component]
	m.init(tp)

	// store it in the cache
	cache[tp] = m
	return m
}

// meta stores meta-information about a specific component
type meta[Component any] struct {
	Name string       // the type name of this component
	Elem reflect.Type // the element type of the component

	CFields map[string]reflect.Type // fields with type C for which C implements component
	IFields map[string]reflect.Type // fields []I where I is an interface that implements component

	DCFields map[string]reflect.Type // fields with type C for which C inside auto field which implement component
	DIFields map[string]reflect.Type // fields []I where I is an interface inside auto field that implements component
}

// init initializes this meta
func (m *meta[Component]) init(tp reflect.Type) {
	var component = reflectx.TypeOf[Component]()

	if tp.Kind() != reflect.Pointer && tp.Elem().Kind() != reflect.Struct {
		panic("GetMeta: Type (" + tp.String() + ") must be backed by a pointer to slice")
	}

	m.Elem = tp.Elem()
	m.Name = nameOf(m.Elem)

	m.CFields = make(map[string]reflect.Type)
	m.IFields = make(map[string]reflect.Type)
	scanForFields(component, m.Name, m.Elem, false, m.CFields, m.IFields)

	// check if we have a dependencies field of struct type
	dependenciesField, ok := m.Elem.FieldByName(dependencies)
	if !ok {
		return
	}

	if dependenciesField.Type.Kind() != reflect.Struct {
		panic("GetMeta: " + dependencies + " field (" + m.Name + ") is not a struct")
	}

	// and initialize the type map of the given map
	m.DCFields = make(map[string]reflect.Type)
	m.DIFields = make(map[string]reflect.Type)
	scanForFields(component, m.Name, dependenciesField.Type, true, m.DCFields, m.DIFields)
}

// scanForFields scans the structtype for fields of component-like fields.
// they are then writen to the cFields and iFields maps.
// inDependenciesStruct indicates if we are inside a dependency struct
func scanForFields(component reflect.Type, elem string, structType reflect.Type, inDependenciesStruct bool, cFields map[string]reflect.Type, iFields map[string]reflect.Type) {
	count := structType.NumField()
	for i := 0; i < count; i++ {
		field := structType.Field(i)

		if !inDependenciesStruct && field.Tag.Get("auto") != "true" {
			continue
		}
		if inDependenciesStruct && field.Tag != "" {
			panic("GetMeta: " + dependencies + " field (" + elem + ") contains field (" + field.Name + ") with tag")
		}

		tp := field.Type
		name := field.Name

		switch {
		case implementsComponent(component, tp):
			cFields[name] = tp
		case implementsSlice(component, tp):
			iFields[name] = tp.Elem()
		case inDependenciesStruct:
			panic("GetMeta: " + dependencies + " field (" + elem + ") contains non-auto fields")
		}
	}
}

func implementsComponent(component reflect.Type, tp reflect.Type) bool {
	return tp.Implements(component) && tp.Kind() == reflect.Pointer && tp.Elem().Kind() == reflect.Struct
}

func implementsSlice(component reflect.Type, tp reflect.Type) bool {
	return tp.Kind() == reflect.Slice && tp.Elem().Kind() == reflect.Interface && tp.Elem().Implements(component)
}

func nameOf(tp reflect.Type) string {
	return tp.PkgPath() + "." + tp.Name()
}

// New creates a new ComponentDescription
func (m meta[Component]) New() Component {
	return reflect.New(m.Elem).Interface().(Component)
}

// NeedsInitComponent
func (m meta[Component]) NeedsInitComponent() bool {
	return len(m.CFields) > 0 || len(m.IFields) > 0 || len(m.DCFields) > 0 || len(m.DIFields) > 0
}

// name of the dependencies field
const dependencies = "Dependencies"

// InitComponent sets up the fields of the given instance of a component.
func (m meta[Component]) InitComponent(instance reflect.Value, all []Component) {
	elem := instance.Elem()
	dependenciesElem := elem.FieldByName(dependencies)

	// assign the component fields
	for field, eType := range m.CFields {
		c := collection.First(all, func(c Component) bool {
			return reflect.TypeOf(c).AssignableTo(eType)
		})

		field := elem.FieldByName(field)
		field.Set(reflect.ValueOf(c))
	}
	for field, eType := range m.DCFields {
		c := collection.First(all, func(c Component) bool {
			return reflect.TypeOf(c).AssignableTo(eType)
		})

		field := dependenciesElem.FieldByName(field)
		field.Set(reflect.ValueOf(c))
	}

	// assign the interface subtypes
	registryR := reflect.ValueOf(all)
	for field, eType := range m.IFields {
		cs := filterSubtype(registryR, eType)
		field := elem.FieldByName(field)
		field.Set(cs)
	}
	for field, eType := range m.DIFields {
		cs := filterSubtype(registryR, eType)
		field := dependenciesElem.FieldByName(field)
		field.Set(cs)
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
		if element.Elem().Type().Implements(iface) {
			result = reflect.Append(result, element.Elem().Convert(iface))
		}
	}
	return result
}
