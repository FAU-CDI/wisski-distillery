// Package component holds the main abstraction for components.
package component

import (
	"reflect"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// Components represents a logical subsystem of the distillery.
// A Component should be implemented as a pointer to a struct.
// Every component must embed [Base] and should be initialized using [Init] inside a [lazy.Pool].
//
// By convention these are defined within their corresponding subpackage.
// This subpackage also contains all required resources.
type Component interface {
	// Name returns the name of this component
	// Name should be implemented by the [ComponentBase] struct.
	Name() string

	// getBase returns the underlying ComponentBase object of this Component.
	// It is used internally during initialization
	getBase() *Base
}

// Base is embedded into every Component
type Base struct {
	name  string // name is the name of this component
	Still        // the underlying still of the distillery
}

//lint:ignore U1000 used to implement the private methods of [Component]
func (cb *Base) getBase() *Base {
	return cb
}

// Init initialzes a new componeont Component with the provided still.
// Init is only initended to be used within a lazy.Pool[Component,Still].
func Init(component Component, core Still) Component {
	base := component.getBase() // pointer to a struct
	base.Still = core
	base.name = strings.ToLower(reflect.TypeOf(component).Elem().Name())
	return component
}

func (cb Base) Name() string {
	return cb.name
}

// Still represents the central part of a distillery.
// It is used inside the main distillery struct, as well as every component via [ComponentBase].
type Still struct {
	Environment environment.Environment // environment to use for reading / writing to and from the distillery
	Config      *config.Config          // the configuration of the distillery
}
