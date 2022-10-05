package component

import (
	"reflect"

	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

// Pool holds a pool of components and provides factilities to create and access them.
// See [Pool.All], [ExportComponents] and [ExportComponent].
type Pool struct {
	all lazy.Lazy[[]Component] // all components
}

// All initializes or returns all components stored in this pool.
//
// The All function should return an array of calls to [Make] with the provided context.
// Multiple calls to the this method return the same return value.
func (p *Pool) All(All func(context *PoolContext) []Component) []Component {
	return p.all.Get(func() []Component {
		// create a new context
		context := &PoolContext{
			all:   All,
			cache: make(map[string]Component),
		}

		// and process them all
		all := context.all(context)
		context.process(all)
		return all
	})
}

// PoolContext is a context used during [Make].
// It should not be initialized by a user.
type PoolContext struct {
	all func(context *PoolContext) []Component // function to return all components

	cache map[string]Component // cached components
	queue []delayedInit        // init queue
}

type delayedInit struct {
	meta  meta
	value reflect.Value
}

func (di delayedInit) Do(all []Component) {
	di.meta.InitComponent(di.value, all)
}

// process processes the queue in this process
func (p *PoolContext) process(all []Component) {
	index := 0
	for len(p.queue) > index {
		p.queue[index].Do(all)
		index++
	}
	p.queue = nil
}

// Make creates or returns a cached component of the given Context.
//
// Components are initialized by first calling the init function.
// Then all component-like fields of fields are filled with their appropriate components.
//
// A component-like field has one of the following types:
//
// - A pointer to a struct type that implements component
// - A slice type of an interface type that implements component
//
// These fields are initialized in an undefined order during initialization.
// The init function may not rely on these existing.
// Furthermore, the init function may not cause other components to be initialized.
//
// The init function may be nil, indicating that no additional initialization is required.
func Make[C Component](context *PoolContext, core Core, init func(component C)) C {
	// get a description of the type
	cd := getMeta[C]()

	// if an instance already exists, return it!
	if instance, ok := context.cache[cd.Name]; ok {
		return instance.(C)
	}

	// make sure that we have an array of components
	if context.cache == nil {
		context.cache = make(map[string]Component)
	}

	// create a fresh (empty) instance
	context.cache[cd.Name] = cd.New()
	instance := context.cache[cd.Name].(C)

	// do the core and self-initialization
	instance.getBase().Core = core

	if init != nil {
		init(instance)
	}

	if cd.NeedsInitComponent() {
		context.queue = append(context.queue, delayedInit{
			meta:  cd,
			value: reflect.ValueOf(instance),
		})
	}

	return instance
}

//
// PUBLIC FUNCTIONS
//

// ExportComponents exports all components that are a C from the pool.
//
// All should be the function of the core that initializes all components.
// All should only make calls to [InitComponent].
func ExportComponents[C Component](p *Pool, All func(context *PoolContext) []Component) []C {
	components := p.All(All)

	results := make([]C, 0, len(components))
	for _, comp := range components {
		if cc, ok := comp.(C); ok {
			results = append(results, cc)
		}
	}
	return results
}

// ExportComponent exports the first component that is a C from the pool.
//
// All should be the function of the core that initializes all components.
// All should only make calls to [InitComponent].
func ExportComponent[C Component](p *Pool, All func(context *PoolContext) []Component) C {
	components := p.All(All)

	for _, comp := range components {
		if cc, ok := comp.(C); ok {
			return cc
		}
	}

	var c C
	return c
}
