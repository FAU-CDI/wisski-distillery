//spellchecker:words actions
package actions

//spellchecker:words context encoding json github wisski distillery internal component auth scopes models
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

type Rebuild struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*Rebuild)(nil)
)

func (*Rebuild) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "rebuild",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 1,
		},
	}
}

func (r *Rebuild) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	// read the flags of the instance to be rebuilt
	var system models.System
	if err := json.Unmarshal([]byte(params[0]), &system); err != nil {
		return nil, fmt.Errorf("failed to unmarshal system properties: %w", err)
	}

	if err := instance.SystemManager().Apply(ctx, out, system); err != nil {
		return nil, fmt.Errorf("failed to apply system properties: %w", err)
	}
	return nil, nil
}
