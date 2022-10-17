// Package component holds the main abstraction for components.
package component

// Component represents a logical subsystem of the distillery.
// Every component must embed [ComponentBase] and should be initialized using [Initialize].
// A Component should be implemented as a pointer to a struct.
//
// By convention these are defined within their corresponding subpackage.
// This subpackage also contains all required resources.
//
// Components are initialized using a [Pool].
type Component interface {
	// Name returns the name of this component.
	// It should correspond to the appropriate subpackage.
	Name() string

	// getBase returns the embedded ComponentBase struct.
	getBase() *ComponentBase
}

// ComponentBase should be embedded into every component
type ComponentBase struct {
	Core // the underlying stillage of the distillery
}

//lint:ignore U1000 used to implement the private methods of [Component]
func (cb *ComponentBase) getBase() *ComponentBase {
	return cb
}
