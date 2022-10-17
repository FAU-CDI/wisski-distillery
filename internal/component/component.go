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
	// Name returns the name of this component
	// Name should be implemented by the [ComponentBase] struct.
	Name() string

	// getComponentBase returns the underlying ComponentBase object of this Component.
	// It is used internally during initialization
	getComponentBase() *ComponentBase
}
