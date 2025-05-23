//spellchecker:words actions
package actions

//spellchecker:words context github wisski distillery internal component auth scopes
import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

// WebsocketAction is a component that can be run via the websocket api.
type WebsocketAction interface {
	component.Component

	Action() Action
	Act(ctx context.Context, in io.Reader, out io.Writer, params ...string) (any, error)
}

// WebsocketInstanceAction is like [WebsocketAction] but is specific to a WissKI instance.
type WebsocketInstanceAction interface {
	component.Component

	Action() InstanceAction
	Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error)
}

// Action represents information about an action.
type Action struct {
	Name string

	Scope      scopes.Scope
	ScopeParam string
	NumParams  int
}

type InstanceAction struct {
	Action
}
