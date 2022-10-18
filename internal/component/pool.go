package component

import (
	"sync"

	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

// Pool holds a pool of components.
type Pool struct {
	pool     lazy.Pool[Component, Still]
	poolInit sync.Once
}

func (pool *Pool) init() {
	pool.poolInit.Do(func() {
		pool.pool.Init = Init
	})
}

type PoolContext = *lazy.PoolContext[Component]
type AllFunc = func(context PoolContext) []Component

func (pool *Pool) All(core Still, init func(context PoolContext) []Component) []Component {
	pool.init()
	return pool.pool.All(core, init)
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
func Make[C Component](context PoolContext, core Still, init func(component C)) C {
	return lazy.Make(context, init)
}

// ExportAll exports all components that are a C from the pool.
//
// All should be the function of the core that initializes all components.
// All should only make calls to [InitComponent].
func ExportAll[C Component](pool *Pool, core Still, All AllFunc) []C {
	pool.init()
	return lazy.ExportComponents[Component, Still, C](&pool.pool, core, All)
}

// Export exports the first component that is a C from the pool.
//
// All should be the function of the core that initializes all components.
// All should only make calls to [InitComponent].
func Export[C Component](pool *Pool, core Still, All AllFunc) C {
	pool.init()
	return lazy.ExportComponent[Component, Still, C](&pool.pool, core, All)
}
