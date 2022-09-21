// Package component holds the main abstraction for components.
package component

import (
	"reflect"

	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

// Component represents a logical subsystem of the distillery.
// Every component must embed [ComponentBase] and should be initialized using [Initialize].
//
// By convention these are defined within their corresponding subpackage.
// This subpackage also contains all required resources.
// Furthermore, a component is typically instantiated using a call on the ["distillery.Distillery"] struct.
//
// For example, the web.Web component lives in the web package and can be created like:
//
//	var dis Distillery
//  web := dis.Web()
type Component interface {
	// Name returns the name of this component.
	// It should correspond to the appropriate subpackage.
	Name() string

	// getBase returns the embedded ComponentBase struct.
	getBase() *ComponentBase
}

// ComponentBase should be embedded into every component
type ComponentBase struct {
	Core // the core of the associated distillery
}

//lint:ignore U1000 used to implement the private methods of [Component]
func (cb *ComponentBase) getBase() *ComponentBase {
	return cb
}

// Initialize makes or returns a component based on a lazy.
//
// C is the type of component to initialize. It must be backed by a pointer, or Initialize will panic.
//
// dis is the distillery to initialize components for
// field is a pointer to the appropriate struct field within the distillery components
// init is called with a new non-nil component to initialize it.
// It may be nil, to indicate no additional initialization is required.
//
// makeComponent returns the new or existing component instance
func Initialize[C Component](core Core, field *lazy.Lazy[C], init func(C)) C {

	// get the typeof C and make sure that it is a pointer type!
	typC := reflect.TypeOf((*C)(nil)).Elem()
	if typC.Kind() != reflect.Pointer {
		panic("Initialize: C must be backed by a pointer")
	}

	// return the field
	return field.Get(func() (c C) {
		c = reflect.New(typC.Elem()).Interface().(C)
		if init != nil {
			init(c)
		}

		base := c.getBase()
		base.Core = core

		return
	})
}
