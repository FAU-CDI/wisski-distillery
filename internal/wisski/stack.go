package wisski

import (
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/meta"
	"github.com/tkw1536/goprogram/stream"
)

//go:embed all:instances/barrel instances/barrel.env
var barrelResources embed.FS

// Barrel returns a stack representing the running WissKI Instance
func (wisski *WissKI) Barrel() component.StackWithResources {
	return component.StackWithResources{
		Stack: component.Stack{
			Dir: wisski.FilesystemBase,
			Env: wisski.Core.Environment,
		},

		Resources:   barrelResources,
		ContextPath: filepath.Join("instances", "barrel"),
		EnvPath:     filepath.Join("instances", "barrel.env"),

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": wisski.Core.Config.DockerNetworkName,

			"SLUG":          wisski.Slug,
			"VIRTUAL_HOST":  wisski.Domain(),
			"HTTPS_ENABLED": wisski.Core.Config.HTTPSEnabledEnv(),

			"DATA_PATH":                   filepath.Join(wisski.FilesystemBase, "data"),
			"RUNTIME_DIR":                 wisski.Core.Config.RuntimeDir(),
			"GLOBAL_AUTHORIZED_KEYS_FILE": wisski.Core.Config.GlobalAuthorizedKeysFile,
		},

		MakeDirs: []string{"data", ".composer"},

		TouchFiles: []string{
			filepath.Join("data", "authorized_keys"),
		},
	}
}

// TODO: Move this to time.Time
var lastRebuild = meta.StorageFor[int64]("lastRebuild")

func (wisski *WissKI) LastRebuild() (t time.Time, err error) {
	epoch, err := lastRebuild(wisski.storage()).Get()
	if err == meta.ErrMetadatumNotSet {
		return t, nil
	}
	if err != nil {
		return t, err
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (wisski *WissKI) setLastRebuild() error {
	return lastRebuild(wisski.storage()).Set(time.Now().Unix())
}

// Build builds or rebuilds the barel connected to this instance.
//
// It also logs the current time into the metadata belonging to this instance.
func (wisski *WissKI) Build(stream stream.IOStream, start bool) error {
	if err := wisski.TryLock(); err != nil {
		return err
	}
	defer wisski.Unlock()

	barrel := wisski.Barrel()

	var context component.InstallationContext

	{
		err := barrel.Install(stream, context)
		if err != nil {
			return err
		}
	}

	{
		err := barrel.Update(stream, start)
		if err != nil {
			return err
		}
	}

	// store the current last rebuild
	return wisski.setLastRebuild()
}

//go:embed all:instances/reserve instances/reserve.env
var reserveResources embed.FS

// Reserve returns a stack representing the reserve instance
func (wisski *WissKI) Reserve() component.StackWithResources {
	return component.StackWithResources{
		Stack: component.Stack{
			Dir: wisski.FilesystemBase,
			Env: wisski.Core.Environment,
		},

		Resources:   reserveResources,
		ContextPath: filepath.Join("instances", "reserve"),
		EnvPath:     filepath.Join("instances", "reserve.env"),

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": wisski.Core.Config.DockerNetworkName,

			"SLUG":          wisski.Slug,
			"VIRTUAL_HOST":  wisski.Domain(),
			"HTTPS_ENABLED": wisski.Core.Config.HTTPSEnabledEnv(),
		},
	}
}

// Shell executes a shell command inside the instance.
func (wisski *WissKI) Shell(io stream.IOStream, argv ...string) (int, error) {
	return wisski.Barrel().Exec(io, "barrel", "/bin/sh", append([]string{"/user_shell.sh"}, argv...)...)
}
