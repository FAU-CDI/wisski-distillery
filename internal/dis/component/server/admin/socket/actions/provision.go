package actions

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

func (p *Provision) Act(ctx context.Context, in io.Reader, out io.Writer, params ...string) error {
	// read the flags of the instance to be provisioned
	var flags provision.Flags
	if err := json.Unmarshal([]byte(params[0]), &flags); err != nil {
		return err
	}

	instance, err := p.dependencies.Provision.Provision(
		out,
		ctx,
		flags,
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "URL:      %s\n", instance.URL().String())
	fmt.Fprintf(out, "Username: %s\n", instance.DrupalUsername)
	fmt.Fprintf(out, "Password: %s\n", instance.DrupalPassword)

	return nil
}
