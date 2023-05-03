package dis

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/collection"
	"github.com/tkw1536/pkglib/lifetime"
)

//
//  ==== init ====
//

func (dis *Distillery) init() {
	dis.lifetimeInit.Do(func() {
		dis.lifetime.Init = component.Init

		lifetime.RegisterGroup[component.Backupable](&dis.lifetime)
		lifetime.RegisterGroup[component.Snapshotable](&dis.lifetime)
		lifetime.RegisterGroup[component.DistilleryFetcher](&dis.lifetime)
		lifetime.RegisterGroup[component.Installable](&dis.lifetime)
		lifetime.RegisterGroup[component.Provisionable](&dis.lifetime)
		lifetime.RegisterGroup[component.Routeable](&dis.lifetime)
		lifetime.RegisterGroup[component.Cronable](&dis.lifetime)
		lifetime.RegisterGroup[component.UserDeleteHook](&dis.lifetime)
		lifetime.RegisterGroup[component.Table](&dis.lifetime)
		lifetime.RegisterGroup[component.Menuable](&dis.lifetime)
		lifetime.RegisterGroup[component.ScopeProvider](&dis.lifetime)
	})
}

//
//  ==== registration ====
//

// manual initializes a component from the provided distillery.
func manual[C component.Component](init func(component C)) initFunc {
	return func(context ctx) component.Component {
		return lifetime.Make(context, init)
	}
}

// use is like r, but does not provided additional initialization
func auto[C component.Component](context ctx) component.Component {
	return lifetime.Make[component.Component, C](context, nil)
}

// register returns all components of the distillery
func (dis *Distillery) register(context ctx) []component.Component {
	return collection.MapSlice(
		dis.allComponents(),
		func(f initFunc) component.Component {
			return f(context)
		},
	)
}

// ctx is a context for component initialization
type ctx = *lifetime.InjectorContext[component.Component]

//
//  ==== export ====
//

// export is a convenience function to export a single component
func export[C component.Component](dis *Distillery) C {
	dis.init()
	return lifetime.ExportComponent[component.Component, component.Still, C](&dis.lifetime, dis.Still, dis.register)
}

func exportAll[C component.Component](dis *Distillery) []C {
	dis.init()
	return lifetime.ExportComponents[component.Component, component.Still, C](&dis.lifetime, dis.Still, dis.register)
}

type initFunc = func(context ctx) component.Component
