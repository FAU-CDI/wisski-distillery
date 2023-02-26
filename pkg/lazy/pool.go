package lazy

import (
	"reflect"

	"github.com/tkw1536/pkglib/reflectx"
)

// Pool represents a pool of laziliy initialized and potentially referencing Component instances.
//
// Component must be an interface type, that should be implemented by various pointers to structs.
// Components may reference each other, even circularly.
//
// Each type of struct is considered a singleton an initialized only once.
//
// See [Pool.All], [ExportComponents] and [ExportComponent].
//
// The zero value is ready to use.
type Pool[Component any, InitParams any] struct {
	// Init is called on every component to be initialized.
	Init func(Component, InitParams) Component

	// Analytics are written on the first retrieval operation on this Pool.
	//
	// Contains all groups and structs that are referenced during initialization.
	// To add extra groups, call RegisterPoolGroup.
	Analytics   PoolAnalytics
	extraGroups []reflect.Type

	all Lazy[[]Component]
}

// RegisterPoolGroup registers the given group type to be added to the pools' analytics.
//
// Only groups not referenced during initialization need to be registered explicitly.
func RegisterPoolGroup[Group any, Component any, InitParams any](p *Pool[Component, InitParams]) {
	p.extraGroups = append(p.extraGroups, reflectx.TypeOf[Group]())
}

// All initializes or returns all components stored in this pool.
//
// The All function should return an array of calls to [Make] with the provided context.
// Multiple calls to the this method return the same return value.
func (p *Pool[Component, InitParams]) All(Params InitParams, All func(context *PoolContext[Component]) []Component) []Component {
	return p.all.Get(func() []Component {
		// create a new context
		context := &PoolContext[Component]{
			all: All,
			init: func(c Component) Component {
				if p.Init == nil {
					return c
				}
				return p.Init(c, Params)
			},
			metaCache: make(map[reflect.Type]meta[Component]),
			cache:     make(map[string]Component),
		}

		// and process them all
		all := context.all(context)
		context.process(all)

		// write out analytics
		context.anal(&p.Analytics, p.extraGroups)
		return all
	})
}

// PoolContext is a context used during [Make].
type PoolContext[Component any] struct {
	init func(Component) Component                         // initializes a new component
	all  func(context *PoolContext[Component]) []Component // initializes all components

	metaCache map[reflect.Type]meta[Component]
	cache     map[string]Component     // cached components
	queue     []delayedInit[Component] // init queue
}

type delayedInit[Component any] struct {
	meta  meta[Component]
	value reflect.Value
}

// Process processes all components in the queue
func (p *PoolContext[Component]) process(all []Component) {
	index := 0
	for len(p.queue) > index {
		p.queue[index].Run(all)
		index++
	}
	p.queue = nil
}

func (di *delayedInit[Component]) Run(all []Component) {
	di.meta.InitComponent(di.value, all)
}

// Make creates or returns a cached component of the given Context.
//
// Components are initialized by first
// Then all component-like fields of fields are filled with their appropriate components.
//
// A component-like field has one of the following types:
//
// - A pointer to a struct type that implements component
// - A slice type of an interface type that implements component
//
// Such component-like fields are only initialized if one of the following conditions are met:
//
// - The field has a tag 'auto' with the value `true`
// - The field lives inside a struct field named `Dependencies`
//
// These fields are initialized in an undefined order during initialization.
// The init function may not rely on these existing.
// Furthermore, the init function may not cause other components to be initialized.
//
// The init function may be nil, indicating that no additional initialization is required.
func Make[Component any, ConcreteComponent any](context *PoolContext[Component], init func(component ConcreteComponent)) ConcreteComponent {
	// get a description of the type
	cd := getMeta[Component, ConcreteComponent](context.metaCache)

	// if an instance already exists, return it!
	if instance, ok := context.cache[cd.Name]; ok {
		return any(instance).(ConcreteComponent)
	}

	// create a fresh new instance and store it in the cache
	context.cache[cd.Name] = context.init(cd.New())
	instance := any(context.cache[cd.Name]).(ConcreteComponent)

	// call the passed init function
	if init != nil {
		init(instance)
	}

	// and queue it up
	if cd.NeedsInitComponent() {
		context.queue = append(context.queue, delayedInit[Component]{
			meta:  cd,
			value: reflect.ValueOf(instance),
		})
	}

	return instance
}

//
// PUBLIC FUNCTIONS
//

// ExportComponents exports all components that are a ConcreteComponentType from the pool.
//
// All should be the function of the core that initializes all components.
// All should only make calls to [InitComponent].
func ExportComponents[Component any, InitParams any, ConcreteComponentType any](
	p *Pool[Component, InitParams],
	Params InitParams,
	All func(context *PoolContext[Component]) []Component,
) []ConcreteComponentType {
	components := p.All(Params, All)

	results := make([]ConcreteComponentType, 0, len(components))
	for _, comp := range components {
		if cc, ok := any(comp).(ConcreteComponentType); ok {
			results = append(results, cc)
		}
	}
	return results
}

// ExportComponent exports the first component that is a ConcreteComponent from the pool.
//
// All should be the function of the core that initializes all components.
// All should only make calls to [InitComponent].
func ExportComponent[Component any, InitParams any, ConcreteComponentType any](
	pool *Pool[Component, InitParams],
	Params InitParams,
	All func(context *PoolContext[Component]) []Component,
) ConcreteComponentType {
	components := pool.All(Params, All)

	for _, comp := range components {
		if cc, ok := any(comp).(ConcreteComponentType); ok {
			return cc
		}
	}

	panic("ExportComponent: Attempted to export unregistered component")
}
