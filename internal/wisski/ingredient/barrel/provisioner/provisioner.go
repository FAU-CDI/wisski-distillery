package provisioner

import (
	"errors"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/alessio/shellescape"
	"github.com/tkw1536/goprogram/stream"
)

// Provisioner provides provisioning for a barrel
// NOTE(twiesing): This should be refactored to not use the provision script.
// Instead, this should code directly defined in go.
type Provisioner struct {
	ingredient.Base
	Barrel *barrel.Barrel
}

// Provision provisions an instance, assuming that the required databases already exist.
func (provision *Provisioner) Provision(io stream.IOStream) error {

	// build the container
	if err := provision.Barrel.Build(io, false); err != nil {
		return err
	}

	provisionParams := []string{
		provision.Domain(),

		provision.SqlDatabase,
		provision.SqlUsername,
		provision.SqlPassword,

		provision.GraphDBRepository,
		provision.GraphDBUsername,
		provision.GraphDBPassword,

		provision.DrupalUsername,
		provision.DrupalPassword,

		"", // TODO: DrupalVersion
		"", // TODO: WissKIVersion
	}

	// escape the parameter
	for i, param := range provisionParams {
		provisionParams[i] = shellescape.Quote(param)
	}

	// figure out the provision script
	// TODO: Move the provision script into the control plane!
	provisionScript := "sudo PATH=$PATH -u www-data /bin/bash /provision_container.sh " + strings.Join(provisionParams, " ")

	code, err := provision.Barrel.Stack().Run(io, true, "barrel", "/bin/bash", "-c", provisionScript)
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New("unable to run provision script")
	}

	return nil
}
