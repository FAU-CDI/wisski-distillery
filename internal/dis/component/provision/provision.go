package provision

import (
	"context"
	"errors"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

type Provision struct {
	component.Base
	Dependencies struct {
		Instances     *instances.Instances
		Provisionable []component.Provisionable
	}
}

// ProvisionFlags are flags for a new instance
type ProvisionFlags struct {
	Slug string
}

var ErrInstanceAlreadyExists = errors.New("instance with provided slug already exists")

// Provision provisions a new docker compose instance.
func (pv *Provision) Provision(progress io.Writer, ctx context.Context, flags ProvisionFlags) (*wisski.WissKI, error) {
	// check that it doesn't already exist
	logging.LogMessage(progress, ctx, "Provisioning new WissKI instance %s", flags.Slug)
	if exists, err := pv.Dependencies.Instances.Has(ctx, flags.Slug); err != nil || exists {
		return nil, ErrInstanceAlreadyExists
	}

	// make it in-memory
	instance, err := pv.Dependencies.Instances.Create(flags.Slug)
	if err != nil {
		return nil, err
	}

	// check that the base directory does not exist
	logging.LogMessage(progress, ctx, "Checking that base directory %s does not exist", instance.FilesystemBase)
	if fsx.IsDirectory(instance.FilesystemBase) {
		return nil, ErrInstanceAlreadyExists
	}

	// Store in the instances table!
	if err := logging.LogOperation(func() error {
		if err := instance.Bookkeeping().Save(ctx); err != nil {
			return err
		}

		return nil
	}, progress, ctx, "Updating bookkeeping database"); err != nil {
		return nil, err
	}

	// create all the resources!
	if err := logging.LogOperation(func() error {
		domain := instance.Domain()
		for _, pc := range pv.Dependencies.Provisionable {
			logging.LogMessage(progress, ctx, "Provisioning %s resources", pc.Name())
			err := pc.Provision(ctx, instance.Instance, domain)
			if err != nil {
				return err
			}
		}

		return nil
	}, progress, ctx, "Provisioning instance-specific resources"); err != nil {
		return nil, err
	}

	// run the provision script
	if err := logging.LogOperation(func() error {
		return instance.Provisioner().Provision(ctx, progress)
	}, progress, ctx, "Running setup scripts"); err != nil {
		return nil, err
	}

	// start the container!
	logging.LogMessage(progress, ctx, "Starting Container")
	if err := instance.Barrel().Stack().Up(ctx, progress); err != nil {
		return nil, err
	}

	// and return the instance
	return instance, nil
}
