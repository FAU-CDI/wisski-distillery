package dis

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/tkw1536/goprogram/lib/collection"
)

//
//  ==== init ====
//

func (dis *Distillery) init() {
	dis.poolInit.Do(func() {
		dis.pool.Init = component.Init
	})
}

//
//  ==== registration ====
//

// manual initializes a component from the provided distillery.
func manual[C component.Component](init func(component C)) initFunc {
	return func(context ctx) component.Component {
		return lazy.Make(context, init)
	}
}

// use is like r, but does not provided additional initialization
func auto[C component.Component](context ctx) component.Component {
	return lazy.Make[component.Component, C](context, nil)
}

// register returns all components of the distillery
func (dis *Distillery) register(context ctx) []component.Component {
	dis.poolInit.Do(func() {
		dis.pool.Init = component.Init
	})

	return collection.MapSlice(
		dis.allComponents(),
		func(f initFunc) component.Component {
			return f(context)
		},
	)
}

// ctx is a context for component initialization
type ctx = *lazy.PoolContext[component.Component]

//
//  ==== export ====
//

// export is a convenience function to export a single component
func export[C component.Component](dis *Distillery) C {
	dis.init()
	return lazy.ExportComponent[component.Component, component.Still, C](&dis.pool, dis.Still, dis.register)
}

func exportAll[C component.Component](dis *Distillery) []C {
	dis.init()
	return lazy.ExportComponents[component.Component, component.Still, C](&dis.pool, dis.Still, dis.register)
}

type initFunc = func(context ctx) component.Component
