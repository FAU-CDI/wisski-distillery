package env

import (
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

const ResolverPrefixFile = "prefix.cfg"

func (dis *Distillery) ResolverStack() stack.Installable {
	stack := dis.asCoreStack(stack.Installable{
		Stack: stack.Stack{
			Name: "resolver",
		},

		EnvFileContext: map[string]string{
			"VIRTUAL_HOST":      dis.DefaultVirtualHost(),
			"LETSENCRYPT_HOST":  dis.DefaultLetsencryptHost(),
			"LETSENCRYPT_EMAIL": dis.Config.CertbotEmail,
			"PREFIX_FILE":       "", // set below!
			"DEFAULT_DOMAIN":    dis.Config.DefaultDomain,
			"LEGACY_DOMAIN":     strings.Join(dis.Config.SelfExtraDomains, ","),
		},

		TouchFiles: []string{ResolverPrefixFile},
	})
	stack.EnvFileContext["PREFIX_FILE"] = filepath.Join(stack.Dir, ResolverPrefixFile)
	return stack
}

func (dis *Distillery) ResolverStackPath() string {
	return dis.ResolverStack().Dir
}

func (dis Distillery) ResolverPrefixConfig() string {
	return filepath.Join(dis.ResolverStackPath(), ResolverPrefixFile)
}
