package docker

import (
	"context"
	"path/filepath"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/loader"
	ctypes "github.com/compose-spec/compose-go/types"
	"golang.org/x/exp/slices"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// ComposeProject represents a compose project
type ComposeProject = *ctypes.Project

// LoadComposeProject loads a docker compose project from the given path.
// The returned name indicates the name, as would be found by the 'docker compose' executable.
// If the project could not be found, an appropriate error is returned.
//
// NOTE: This intentionally omits using any kind of api for docker compose.
// This saves a *a lot* of dependencies./
func (docker *Docker) LoadComposeProject(path string) (project ComposeProject, err error) {
	ppath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	opts, err := cli.NewProjectOptions(
		/* configs = */ nil,
		cli.WithWorkingDirectory(ppath),
		cli.WithDefaultConfigPath,
		cli.WithName(loader.NormalizeProjectName(filepath.Base(ppath))),
		cli.WithDotEnv,
	)
	if err != nil {
		return nil, err
	}

	proj, err := cli.ProjectFromOptions(opts)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

// Containers loads the compose project at path, connects to the docker daemon, and then lists all containers belonging to the given services.
// If services is empty, all containers belonging to any service are returned.
func (docker *Docker) Containers(ctx context.Context, path string, services ...string) (containers []types.Container, err error) {
	proj, err := docker.LoadComposeProject(path)
	if err != nil {
		return nil, err
	}

	client, err := docker.APIClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return docker.containers(ctx, proj, client, false, services...)
}

const (
	projectLabel    = "com.docker.compose.project"
	workingDirLabel = "com.docker.compose.project.working_dir"
	serviceLabel    = "com.docker.compose.service"
)

// containers uses the given project and client to find containers belonging to the provided services.
// If all is false, only running containers are returned.
// If all is true, all containers are returned.
//
// services optionally filters the returned containers by the services they belong to.
// If services is empty, all containers are returned, else containers belonging to any of the services included.
func (*Docker) containers(ctx context.Context, project ComposeProject, client DockerClient, all bool, services ...string) ([]types.Container, error) {
	// build filters
	f := filters.NewArgs(
		filters.Arg("label", projectLabel+"="+project.Name),
		filters.Arg("label", workingDirLabel+"="+project.WorkingDir),
	)

	// if there is only one label requested, filter it in the query!
	if len(services) == 1 {
		f.Add("label", serviceLabel+"="+services[0])
	}

	// find the containers
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		All:     all,
		Filters: f,
	})
	if err != nil {
		return nil, err
	}

	// for all services or exactly one service (case above)
	// we can immediatly return!
	if len(services) <= 1 {
		return containers, err
	}

	// make a map of services that were requested
	req := make(map[string]struct{}, len(services))
	for _, service := range services {
		req[service] = struct{}{}
	}

	// filter the containers by what was requested
	result := containers[:0]
	for _, container := range containers {
		service, ok := container.Labels[serviceLabel]
		if !ok {
			continue
		}
		if _, ok := req[service]; ok {
			result = append(result, container)
		}
	}

	// and clip the results
	return slices.Clip(result), nil
}
