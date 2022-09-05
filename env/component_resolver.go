package env

import (
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

// ResolverComponent represents the 'resolver' layer belonging to a distillery
type ResolverComponent struct {
	ConfigName string // Filename of the configuration file

	dis *Distillery
}

// Resolver returns the ResolverComponent belonging to this distillery
func (dis *Distillery) Resolver() ResolverComponent {
	return ResolverComponent{
		ConfigName: "prefix.cfg",

		dis: dis,
	}
}

func (ResolverComponent) Name() string {
	return "resolver"
}

func (resolver ResolverComponent) Stack() stack.Installable {
	stack := resolver.dis.makeComponentStack(resolver, stack.Installable{
		EnvFileContext: map[string]string{
			"VIRTUAL_HOST":      resolver.dis.DefaultVirtualHost(),
			"LETSENCRYPT_HOST":  resolver.dis.DefaultLetsencryptHost(),
			"LETSENCRYPT_EMAIL": resolver.dis.Config.CertbotEmail,
			"PREFIX_FILE":       "", // set below!
			"DEFAULT_DOMAIN":    resolver.dis.Config.DefaultDomain,
			"LEGACY_DOMAIN":     strings.Join(resolver.dis.Config.SelfExtraDomains, ","),
		},

		TouchFiles: []string{resolver.ConfigName},
	})
	stack.EnvFileContext["PREFIX_FILE"] = filepath.Join(stack.Dir, resolver.ConfigName)
	return stack
}

func (resolver ResolverComponent) Path() string {
	return resolver.Stack().Dir
}

func (resolver ResolverComponent) ConfigPath() string {
	return filepath.Join(resolver.Path(), resolver.ConfigName)
}
