package component

import (
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/goprogram/stream"
)

// Installable implements an installable component.
type Installable interface {
	Component

	// Stack can be used to gain access to the "docker compose" stack.
	//
	// This should internally call [ComponentBase.MakeStack]
	Stack(env environment.Environment) StackWithResources

	// Context returns a new InstallationContext to be used during installation from the command line.
	// Typically this should just pass through the parent, but might perform other tasks.
	Context(parent InstallationContext) InstallationContext
}

// Updatable represents a component with an Update method.
type Updatable interface {
	Component

	// Update updates or initializes the provided components.
	// It is called after the component has been installed (if applicable).
	//
	// It may send output to the provided stream.
	//
	// Updating should be idempotent, meaning running it multiple times must not break the existing system.
	Update(stream stream.IOStream) error
}
