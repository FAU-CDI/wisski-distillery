package env

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/embed"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

// TODO: Move everything into specific subpackages

// Stacks returns the Stacks of this distillery
func (dis *Distillery) Components() []Component {
	// TODO: Do we want to cache these components?
	return []Component{
		dis.Web(),
		dis.Self(),
		dis.Resolver(),
		dis.Dis(),
		dis.SSH(),
		dis.Triplestore(),
		dis.SQL(),
	}
}

// Component represents a component of the distillery
type Component interface {
	Name() string // Name is the name of this component

	Stack() stack.Installable                                           // Stack returns the installable stack representing this component
	Context(parent stack.InstallationContext) stack.InstallationContext // context for installation

	Path() string // Path returns the path to this component
}

// asCoreStack treats the provided stack as a core component of this distillery.
func (dis *Distillery) makeComponentStack(component Component, stack stack.Installable) stack.Installable {
	stack.Dir = dis.getComponentPath(component)

	name := component.Name()

	// TODO: This writes out resources.
	// Should migrate this directly!
	if stack.Resources == nil {
		stack.Resources = embed.ResourceEmbed
		stack.ContextPath = filepath.Join("resources", "compose", name)
		stack.EnvPath = filepath.Join("resources", "templates", "docker-env", name)
	}

	return stack
}

func (dis *Distillery) getComponentPath(component Component) string {
	return filepath.Join(dis.Config.DeployRoot, "core", component.Name())
}
