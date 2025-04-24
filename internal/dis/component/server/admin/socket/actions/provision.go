//spellchecker:words actions
package actions

//spellchecker:words context encoding json github wisski distillery internal component auth scopes provision
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
)

type Provision struct {
	component.Base
	dependencies struct {
		Provision *provision.Provision
	}
}

var (
	_ WebsocketAction = (*Provision)(nil)
)

func (*Provision) Action() Action {
	return Action{
		Name:      "provision",
		Scope:     scopes.ScopeUserAdmin,
		NumParams: 1,
	}
}

type ProvisionResult struct {
	URL            string
	DrupalUsername string
	DrupalPassword string
}

func (p *Provision) Act(ctx context.Context, in io.Reader, out io.Writer, params ...string) (any, error) {
	// read the flags of the instance to be provisioned
	var flags provision.Flags
	if err := json.Unmarshal([]byte(params[0]), &flags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provision flags: %w", err)
	}

	instance, err := p.dependencies.Provision.Provision(
		out,
		ctx,
		flags,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to provision system: %w", err)
	}

	result := ProvisionResult{
		URL:            instance.URL().String(),
		DrupalUsername: instance.DrupalUsername,
		DrupalPassword: instance.DrupalPassword,
	}

	fmt.Fprintf(out, "URL:      %s\n", result.URL)
	fmt.Fprintf(out, "Username: %s\n", result.DrupalUsername)
	fmt.Fprintf(out, "Password: %s\n", result.DrupalPassword)

	return result, nil
}
