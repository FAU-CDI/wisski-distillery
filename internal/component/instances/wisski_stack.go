package instances

import (
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/tkw1536/goprogram/stream"
)

//go:embed all:instances/barrel instances/barrel.env
var barrelResources embed.FS

// Barrel returns a stack representing the running WissKI Instance
func (wisski *WissKI) Barrel() component.StackWithResources {
	return component.StackWithResources{
		Stack: component.Stack{
			Dir: wisski.FilesystemBase,
			Env: wisski.instances.Environment,
		},

		Resources:   barrelResources,
		ContextPath: filepath.Join("instances", "barrel"),
		EnvPath:     filepath.Join("instances", "barrel.env"),

		EnvContext: map[string]string{
			"DATA_PATH": filepath.Join(wisski.FilesystemBase, "data"),

			"SLUG":         wisski.Slug,
			"VIRTUAL_HOST": wisski.Domain(),

			"LETSENCRYPT_HOST":  wisski.instances.Config.IfHttps(wisski.Domain()),
			"LETSENCRYPT_EMAIL": wisski.instances.Config.IfHttps(wisski.instances.Config.CertbotEmail),

			"RUNTIME_DIR":                 wisski.instances.Config.RuntimeDir(),
			"GLOBAL_AUTHORIZED_KEYS_FILE": wisski.instances.Config.GlobalAuthorizedKeysFile,
		},

		MakeDirs: []string{"data", ".composer"},

		TouchFiles: []string{
			filepath.Join("data", "authorized_keys"),
		},
	}
}

const KeyLastRebuild MetaKey = "lastRebuild"

func (wisski *WissKI) LastRebuild() (t time.Time, err error) {
	var epoch int64

	// read the epoch!
	err = wisski.Metadata().Get(KeyLastRebuild, &epoch)
	if err == ErrMetadatumNotSet {
		return t, nil
	}
	if err != nil {
		return t, err
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (wisski *WissKI) setLastRebuild() error {
	return wisski.Metadata().Set(KeyLastRebuild, time.Now().Unix())
}

// Build builds or rebuilds the barel connected to this instance.
//
// It also logs the current time into the metadata belonging to this instance.
func (wisski *WissKI) Build(stream stream.IOStream, start bool) error {
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
			Env: wisski.instances.Environment,
		},

		Resources:   reserveResources,
		ContextPath: filepath.Join("instances", "reserve"),
		EnvPath:     filepath.Join("instances", "reserve.env"),

		EnvContext: map[string]string{
			"VIRTUAL_HOST": wisski.Domain(),

			"LETSENCRYPT_HOST":  wisski.instances.Config.IfHttps(wisski.Domain()),
			"LETSENCRYPT_EMAIL": wisski.instances.Config.IfHttps(wisski.instances.Config.CertbotEmail),
		},
	}
}
