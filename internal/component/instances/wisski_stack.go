package instances

import (
	"embed"
	"io/fs"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
)

//go:embed all:instances/barrel instances/barrel.env
var barrelResources embed.FS

// Barrel returns a stack representing the running WissKI Instance
func (wisski WissKI) Barrel() component.StackWithResources {
	return component.StackWithResources{
		Stack: component.Stack{
			Dir: wisski.FilesystemBase,
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

		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs:     []string{"data", ".composer"},

		TouchFiles: []string{
			filepath.Join("data", "authorized_keys"),
		},
	}
}

//go:embed all:instances/reserve instances/reserve.env
var reserveResources embed.FS

// Reserve returns a stack representing the reserve instance
func (wisski WissKI) Reserve() component.StackWithResources {
	return component.StackWithResources{
		Stack: component.Stack{
			Dir: wisski.FilesystemBase,
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
