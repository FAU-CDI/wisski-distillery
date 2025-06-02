//spellchecker:words dockerx
package dockerx

//spellchecker:words context errors slices github wisski distillery compose execx docker types container filters pkglib errorsx stream
import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/FAU-CDI/wisski-distillery/pkg/compose"
	"github.com/FAU-CDI/wisski-distillery/pkg/execx"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/stream"
)

// Stack represents a 'docker compose' stack in the provided directory.
//
// NOTE(twiesing): In the current implementation this requires a 'docker' executable on the system.
// This executable must be capable of the 'docker compose' command.
// In the future the idea is to replace this with a native docker compose client.
//
//nolint:recvcheck
type Stack struct {
	Dir string // Directory this Stack is located in

	Client     *Client // nattive docker client
	Executable string  // Path to the native docker executable to use
}

// Close closes any resources associated with this stack.
func (stack *Stack) Close() error {
	if err := stack.Client.Close(); err != nil {
		return fmt.Errorf("failed to close docker client: %w", err)
	}

	return nil
}

// Project returns the underlying compose project.
func (stack *Stack) Project(ctx context.Context) (compose.Project, error) {
	proj, err := compose.Open(ctx, stack.Dir)
	if err != nil {
		return nil, fmt.Errorf("compose.Open(%q) failed: %w", stack.Dir, err)
	}
	return proj, nil
}

const (
	projectLabel    = "com.docker.compose.project"
	workingDirLabel = "com.docker.compose.project.working_dir"
	serviceLabel    = "com.docker.compose.service"
)

// Containers lists all containers belonging to the given services.
// includeStoppedContainers indicates if non-running contains should be included.
// services optionally filters by service name.
func (stack *Stack) Containers(ctx context.Context, includeStoppedContainers bool, services ...string) ([]container.Summary, error) {
	project, err := stack.Project(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open project: %w", err)
	}

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
	containers, err := stack.Client.ContainerList(ctx, container.ListOptions{
		All:     includeStoppedContainers,
		Filters: f,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	// for all services or exactly one service (case above)
	// we can immediatly return!
	if len(services) <= 1 {
		return containers, nil
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

// Kill kills containers belonging to the given service.
func (ds *Stack) Kill(ctx context.Context, progress io.Writer, service string, signal os.Signal) error {
	containers, err := ds.Containers(ctx, false, service)
	if err != nil {
		return fmt.Errorf("failed to get containers for service %q: %w", service, err)
	}

	// kill each of the containers
	errors := make([]error, len(containers))
	for i, container := range containers {
		err := ds.Client.ContainerKill(ctx, container.ID, signal.String())
		if err != nil {
			errors[i] = fmt.Errorf("failed to kill container %q: %w", container.ID, err)
		}
	}

	return errorsx.Combine(errors...)
}

var errStackUpdatePull = errors.New("Stack.Update: Pull returned non-zero exit code")
var errStackUpdateBuild = errors.New("Stack.Update: Build returned non-zero exit code")

// Update pulls, builds, and then optionally starts this stack.
// This does not have a direct 'docker compose' shell equivalent.
//
// See also Up.
func (ds *Stack) Update(ctx context.Context, progress io.Writer, start bool) error {
	if code := ds.compose(ctx, stream.NonInteractive(progress), "pull")(); code != 0 {
		return errStackUpdatePull
	}

	if code := ds.compose(ctx, stream.NonInteractive(progress), "build", "--pull")(); code != 0 {
		return errStackUpdateBuild
	}

	if start {
		return ds.Up(ctx, progress)
	}
	return nil
}

var (
	errStackAttach                   = errors.New("Stack.Attach: Attach returned non-zero exit code")
	errStackAttachNoRunningContainer = errors.New("no running containers")
)

// Attach attaches to the standard output (and optionally input streams) and redirects them to io until context is closed.
// When multiple running containers exist, picks the first one.
func (ds *Stack) Attach(ctx context.Context, io stream.IOStream, interactive bool, services ...string) error {
	containers, err := ds.Containers(ctx, false, services...)
	if err != nil {
		return fmt.Errorf("failed to get containers: %w", err)
	}

	runningIndex := slices.IndexFunc(containers, func(c container.Summary) bool { return c.State == "running" })
	if runningIndex < 0 {
		return errStackAttachNoRunningContainer
	}

	// build the attach command line
	attachContainerCmd := []string{"attach", "--sig-proxy=false"}
	if !interactive {
		io = io.NonInteractive()
		attachContainerCmd = append(attachContainerCmd, "--no-stdin")
	}
	attachContainerCmd = append(attachContainerCmd, containers[runningIndex].ID)

	// run the command!
	code := ds.docker(ctx, io, attachContainerCmd...)()
	if err := ctx.Err(); err != nil {
		// if the context was closed, then return code 0
		code = 0
	}
	if code != 0 {
		return fmt.Errorf("%w: %d", errStackAttach, code)
	}
	return nil
}

var errStackUp = errors.New("Stack.Up: Up returned non-zero exit code")

// Up creates and starts the containers in this Stack.
// It is equivalent to 'docker compose up --force-recreate --remove-orphans --detach' on the shell.
func (ds Stack) Up(ctx context.Context, progress io.Writer) error {
	if code := ds.compose(ctx, stream.NonInteractive(progress), "up", "--force-recreate", "--remove-orphans", "--detach")(); code != 0 {
		return fmt.Errorf("%w: %d", errStackUp, code)
	}
	return nil
}

type ExecOptions struct {
	Service string
	User    string

	Cmd  string
	Args []string
}

// Exec executes an executable in the provided running service.
// It is equivalent to 'docker compose exec $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds *Stack) Exec(ctx context.Context, io stream.IOStream, options ExecOptions) func() int {
	compose := []string{"exec"}
	if !io.StdinIsATerminal() {
		compose = append(compose, "--no-TTY")
	}

	if options.User != "" {
		compose = append(compose, "--user", options.User)
	}

	compose = append(compose, options.Service)
	compose = append(compose, options.Cmd)
	compose = append(compose, options.Args...)

	return ds.compose(ctx, io, compose...)
}

type RunFlags struct {
	AutoRemove bool
	Detach     bool
}

// Run runs a command in a running container with the given executable.
// It is equivalent to 'docker compose run [--rm] $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds *Stack) Run(ctx context.Context, io stream.IOStream, flags RunFlags, service, command string, args ...string) (int, error) {
	compose := []string{"run"}
	if flags.AutoRemove {
		compose = append(compose, "--rm")
	}
	if !io.StdinIsATerminal() {
		compose = append(compose, "--no-TTY")
	}
	if flags.Detach {
		compose = append(compose, "--detach")
	}

	compose = append(compose, service, command)
	compose = append(compose, args...)

	code := ds.compose(ctx, io, compose...)()
	return code, nil
}

var errStackRestart = errors.New("Stack.Restart: Restart returned non-zero exit code")

// Restart restarts all containers in this Stack.
// It is equivalent to 'docker compose restart' on the shell.
func (ds Stack) Restart(ctx context.Context, progress io.Writer) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "restart")()
	if code != 0 {
		return errStackRestart
	}
	return nil
}

var errStackDown = errors.New("Stack.Down: Down returned non-zero exit code")

// Down stops and removes all containers in this Stack.
// It is equivalent to 'docker compose down -v' on the shell.
func (ds *Stack) Down(ctx context.Context, progress io.Writer) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "down", "-v")()
	if code != 0 {
		return errStackDown
	}
	return nil
}

// DownAll stops and removes all containers in this Stack, and those not defined in the compose file.
// It is equivalent to 'docker compose down -v --remove-orphans' on the shell.
func (ds *Stack) DownAll(ctx context.Context, progress io.Writer) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "down", "-v", "--remove-orphans")()
	if code != 0 {
		return errStackDown
	}
	return nil
}

// compose executes a 'docker compose' command on this stack.
func (ds *Stack) compose(ctx context.Context, io stream.IOStream, args ...string) func() int {
	return ds.docker(ctx, io, append([]string{"compose"}, args...)...)
}

// docker executes a 'docker' command in the directory of this stack.
func (ds *Stack) docker(ctx context.Context, io stream.IOStream, args ...string) func() int {
	if ds.Executable == "" {
		var err error
		ds.Executable, err = execx.LookPathAbs("docker")
		if err != nil {
			return execx.CommandErrorFunc
		}
	}
	return execx.Exec(ctx, io, ds.Dir, ds.Executable, args...)
}
