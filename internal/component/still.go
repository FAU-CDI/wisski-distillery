package component

import (
	"reflect"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// ComponentBase is embedded into every Component
type ComponentBase struct {
	name  string // name is the name of this component
	Still        // the underlying still of the distillery
}

//lint:ignore U1000 used to implement the private methods of [Component]
func (cb *ComponentBase) getComponentBase() *ComponentBase {
	return cb
}

func (cb ComponentBase) Name() string {
	return cb.name
}

// nameOf returns the name of the given component
func nameOf(component Component) string {
	return strings.ToLower(reflect.TypeOf(component).Elem().Name())
}

// Still represents the central part of a distillery.
// It is used inside the main distillery struct, as well as every component via [ComponentBase].
type Still struct {
	Environment environment.Environment // environment to use for reading / writing to and from the distillery
	Config      *config.Config          // the configuration of the distillery
}
