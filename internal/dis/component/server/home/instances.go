package home

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"golang.org/x/sync/errgroup"
)

// loadInstances loads all the instances into the home route
func (home *Home) loadInstances(ctx context.Context) ([]status.WissKI, error) {
	// find all the WissKIs
	wissKIs, err := home.Dependencies.Instances.All(ctx)
	if err != nil {
		return nil, err
	}

	instances := make([]status.WissKI, len(wissKIs))

	// determine their infos
	var eg errgroup.Group
	for i, instance := range wissKIs {
		i := i
		wissKI := instance
		eg.Go(func() (err error) {
			instances[i], err = wissKI.Info().Information(ctx, false)
			return
		})
	}
	eg.Wait()

	// and return the new instances
	return instances, nil
}

// UpdateInstanceList updates the instances list of the home struct
type UpdateInstanceList struct {
	component.Base
	Dependencies struct {
		Home *Home
	}
}

var (
	_ component.Cronable = (*UpdateInstanceList)(nil)
)

func (*UpdateInstanceList) TaskName() string {
	return "instance list"
}

func (ul *UpdateInstanceList) Cron(ctx context.Context) error {
	names, err := ul.Dependencies.Home.instanceMap(ctx)
	if err != nil {
		return err
	}

	ul.Dependencies.Home.instanceNames.Set(names)
	return nil
}

type UpdateHome struct {
	component.Base
	Dependencies struct {
		Home *Home
	}
}

var (
	_ component.Cronable = (*UpdateHome)(nil)
)

func (ur *UpdateHome) TaskName() string {
	return "home instances fetch"
}

func (ur *UpdateHome) Cron(ctx context.Context) error {
	instances, err := ur.Dependencies.Home.loadInstances(ctx)
	if err != nil {
		return err
	}

	ur.Dependencies.Home.homeInstances.Set(instances)
	return nil
}
