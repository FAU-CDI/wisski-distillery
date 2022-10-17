package component

import (
	"sync"

	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

type ComponentPool struct {
	pool     lazy.Pool[Component, Core]
	poolInit sync.Once
}

func (pool *ComponentPool) init() {
	pool.poolInit.Do(func() {
		pool.pool.Init = func(component Component, core Core) Component {
			component.getBase().Core = core
			return component
		}
	})
}

type ComponentPoolContext = *lazy.PoolContext[Component]
type ComponentAllFunc = func(context ComponentPoolContext) []Component

func (pool *ComponentPool) All(core Core, init func(context ComponentPoolContext) []Component) []Component {
	pool.init()
	return pool.pool.All(core, init)
}

// MakeComponent creates or returns a cached component of the given Context.
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
func MakeComponent[C Component](context ComponentPoolContext, core Core, init func(component C)) C {
	return lazy.Make(context, init)
}

// ExportComponents exports all components that are a C from the pool.
//
// All should be the function of the core that initializes all components.
// All should only make calls to [InitComponent].
func ExportComponents[C Component](pool *ComponentPool, core Core, All ComponentAllFunc) []C {
	pool.init()
	return lazy.ExportComponents[Component, Core, C](&pool.pool, core, All)
}

// ExportComponent exports the first component that is a C from the pool.
//
// All should be the function of the core that initializes all components.
// All should only make calls to [InitComponent].
func ExportComponent[C Component](pool *ComponentPool, core Core, All ComponentAllFunc) C {
	pool.init()
	return lazy.ExportComponent[Component, Core, C](&pool.pool, core, All)
}
