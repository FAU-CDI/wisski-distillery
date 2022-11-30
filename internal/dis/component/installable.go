package component

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// Installable implements an installable component.
type Installable interface {
	Component

	// Path returns the path this component is installed at.
	// By convention it is /var/www/deploy/internal/core/${Name()}
	Path() string

	// Stack can be used to gain access to the "docker compose" stack.
	//
	// This should internally call [ComponentBase.MakeStack]
	Stack(env environment.Environment) StackWithResources

	// Context returns a new InstallationContext to be used during installation from the command line.
	// Typically this should just pass through the parent, but might perform other tasks.
	Context(parent InstallationContext) InstallationContext
}

// MakeStack registers the Installable as a stack
func MakeStack(component Installable, env environment.Environment, stack StackWithResources) StackWithResources {
	stack.Env = env
	stack.Dir = component.Path()
	return stack
}

// Updatable represents a component with an Update method.
type Updatable interface {
	Component

	// Update updates or initializes the provided components.
	// It is called after the component has been installed (if applicable).
	//
	// It may send progress to the provided stream.
	//
	// Updating should be idempotent, meaning running it multiple times must not break the existing system.
	Update(ctx context.Context, progress io.Writer) error
}
