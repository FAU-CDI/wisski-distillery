//spellchecker:words component
package component

//spellchecker:words context github wisski distillery dockerx
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
)

// Installable implements an installable component.
type Installable interface {
	Component

	// Path returns the path this component is installed at.
	// By convention it is /var/www/deploy/internal/core/${Name()}
	Path() string

	// OpenStack can be used to gain access to the "docker compose" stack.
	//
	// This should internally call [ComponentBase.MakeStack]
	OpenStack() (StackWithResources, error)

	// Context returns a new InstallationContext to be used during installation from the command line.
	// Typically this should just pass through the parent, but might perform other tasks.
	Context(parent InstallationContext) InstallationContext
}

// OpenStack can be used to implement stack as an installable.
func OpenStack(component Installable, factory dockerx.Factory, stack StackWithResources) (StackWithResources, error) {
	dockerStack, err := dockerx.NewStack(factory, component.Path())
	if err != nil {
		return StackWithResources{}, fmt.Errorf("failed to create stack: %w", err)
	}
	stack.Stack = dockerStack
	return stack, nil
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
