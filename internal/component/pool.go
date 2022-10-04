package component

import (
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/FAU-CDI/wisski-distillery/pkg/rlock"
	"github.com/tkw1536/goprogram/lib/reflectx"
)

// Pool represents a pool of components
type Pool struct {
	rLock rlock.RLock

	// the actual queue of initi functions!
	nested uint64 // is the q active?
	queue  []func(thread int32)

	// global initalization!
	initOnce sync.Once

	// components and lock!
	cLock      sync.Mutex
	components map[string]Component
}

func (p *Pool) init() {
	p.initOnce.Do(func() {
		p.components = make(map[string]Component)
	})
}

// Find finds all components of the specific subtype
func Find[C Component](components []Component) C {
	for _, c := range components {
		if cc, ok := c.(C); ok {
			return cc
		}
	}
	panic("FindComponent: Invalid arguments")
}

// InitComponent initializes a specific component and caches it within the given pool.
//
// Concurrent calls of InitComponent must use a distinct thread parameter.
// Nested calls of InitComponent should use the same thread parameter.
//
// Init may initialize components, but not call functions on them!
func InitComponent[C Component](p *Pool, thread int32, core Core, init func(component C, thread int32)) C {
	p.init()

	p.rLock.Lock(int(thread))
	defer p.rLock.Unlock()

	// get a description of the type
	cd := getComponent[C]()

	// find a field to put the component into
	instance, created := func() (C, bool) {
		p.cLock.Lock()
		defer p.cLock.Unlock()

		// create the component
		field, ok := p.components[cd.Name]
		if ok {
			return field.(C), false
		}

		// create a new component
		p.components[cd.Name] = cd.New().(Component)
		return p.components[cd.Name].(C), true
	}()

	// if we already created the instance, then there is nothing to do
	// as someone else will init it!
	if !created {
		return instance
	}

	// setup the core initialization now!
	instance.getBase().Core = core

	if init == nil {
		return instance
	}

	// if we are in nested mode, then delay the init!
	if !atomic.CompareAndSwapUint64(&p.nested, 0, 1) {
		func() {
			p.queue = append(p.queue, func(thread int32) {
				init(instance, thread)
			})
		}()
		return instance
	}
	defer atomic.StoreUint64(&p.nested, 0)

	// init ourselves first (everything below will be nested)
	init(instance, thread)

	// do all the delayed initializations
	index := 0
	for len(p.queue) > index {
		p.queue[index](thread)
		index++
	}
	p.queue = nil

	// and return the instance
	return instance
}

// getComponent gets the component belonging to a component type
func getComponent[C Component]() (desc component) {
	tp := reflectx.TypeOf[C]()
	if tp.Kind() != reflect.Pointer {
		panic("getComponent: C must be backed by a pointer")
		// should never be reached!
	}
	desc.Elem = tp.Elem()
	desc.Name = desc.Elem.PkgPath() + "." + desc.Elem.Name()
	return
}

// component represents a component
type component struct {
	Elem reflect.Type // the element type of the component
	Name string       // the name of the component
}

// New creates a new ComponentDescription
func (cd component) New() any {
	return reflect.New(cd.Elem).Interface()
}
