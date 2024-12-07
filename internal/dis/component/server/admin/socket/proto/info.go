//spellchecker:words proto
package proto

//spellchecker:words context github wisski distillery internal component auth scopes
import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
)

// Action is anything that can be retrieved from a
type Action struct {
	// NumPara
	NumParams int

	// Scope and ScopeParam indicate the scope required by the caller.
	// TODO(twiesing): Once we actually include scopes, make them dynamic
	Scope      component.Scope
	ScopeParam string

	// Handle handles this action.
	//
	// ctx is closed once the underlying connection is closed.
	// out is an io.Writer that is automatically sent to the client.
	// params holds exactly NumParams parameters.
	Handle func(ctx context.Context, in io.Reader, out io.Writer, params ...string) error
}

// scope returns the actual scope required by this action.
// If the caller did not provide an actual scope, uses ScopeNever
func (action Action) scope() component.Scope {
	if action.Scope == "" {
		return scopes.ScopeNever
	}
	return action.Scope
}
