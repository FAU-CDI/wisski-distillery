package env

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

// TODO: Move everything into specific subpackages

// Stacks returns the Stacks of this distillery
func (dis *Distillery) Components() []Component {
	// TODO: Do we want to cache these stacks?
	return []Component{
		dis.Web(),
		dis.Self(),
		dis.Resolver(),
		dis.SSH(),
		dis.Triplestore(),
		dis.SQL(),
	}
}

// Component represents a component of the distillery
type Component interface {
	Name() string // Name is the name of this component

	Stack() stack.Installable // Stack returns the installable stack representing this component

	Path() string // Path returns the path to this component
}

// asCoreStack treats the provided stack as a core component of this distillery.
func (dis *Distillery) makeComponentStack(component Component, stack stack.Installable) stack.Installable {
	name := component.Name()

	stack.Dir = filepath.Join(dis.Config.DeployRoot, "core", name)

	stack.ContextResource = filepath.Join("resources", "compose", name)
	stack.EnvFileResource = filepath.Join("resources", "templates", "docker-env", name)

	return stack
}
