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
	cd := GetMeta[C]()

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

// GetMeta gets the component belonging to a component type
func GetMeta[C Component]() (meta Meta) {
	tp := reflectx.TypeOf[C]()
	if tp.Kind() != reflect.Pointer {
		panic("GetMeta: C must be backed by a pointer (" + tp.String() + ")")
	}
	meta.Elem = tp.Elem()
	meta.Name = meta.Elem.PkgPath() + "." + meta.Elem.Name()
	return
}

// Meta represents meta information about a component
type Meta struct {
	Elem reflect.Type // the element type of the component
	Name string       // the name of the component
}

// New creates a new ComponentDescription
func (cd Meta) New() any {
	return reflect.New(cd.Elem).Interface()
}
