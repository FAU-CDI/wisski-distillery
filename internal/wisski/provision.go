package wisski

import (
	"errors"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/goprogram/stream"
)

// Provision provisions an instance, assuming that the required databases already exist.
func (wisski *WissKI) Provision(io stream.IOStream) error {

	// build the container
	if err := wisski.Build(io, false); err != nil {
		return err
	}

	provisionParams := []string{
		wisski.Domain(),

		wisski.SqlDatabase,
		wisski.SqlUsername,
		wisski.SqlPassword,

		wisski.GraphDBRepository,
		wisski.GraphDBUsername,
		wisski.GraphDBPassword,

		wisski.DrupalUsername,
		wisski.DrupalPassword,

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

	code, err := wisski.Barrel().Run(io, true, "barrel", "/bin/bash", "-c", provisionScript)
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New("unable to run provision script")
	}

	return nil
}
