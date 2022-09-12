package resolver

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/FAU-CDI/wdresolve"
	"github.com/FAU-CDI/wdresolve/resolvers"
	"github.com/FAU-CDI/wisski-distillery/component"
	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/tkw1536/goprogram/stream"
)

// TODO: Add a 'self-server' concept!

type Resolver struct {
	component.ComponentBase

	ConfigName string // the name to the config file
	Executable string // path to the current executable
}

func (Resolver) Name() string {
	return "resolver"
}

func (resolver Resolver) ConfigPath() string {
	return filepath.Join(resolver.Dir, resolver.ConfigName)
}

//go:embed all:stack resolver.env
var resources embed.FS

func (resolver Resolver) Stack() component.Installable {
	return resolver.ComponentBase.MakeStack(component.Installable{
		Resources:   resources,
		ContextPath: "stack",
		EnvPath:     "resolver.env",

		EnvContext: map[string]string{
			"VIRTUAL_HOST":      resolver.Config.DefaultHost(),
			"LETSENCRYPT_HOST":  resolver.Config.DefaultSSLHost(),
			"LETSENCRYPT_EMAIL": resolver.Config.CertbotEmail,

			"CONFIG_PATH": resolver.Config.ConfigPath,
			"DEPLOY_ROOT": resolver.Config.DeployRoot,

			"GLOBAL_AUTHORIZED_KEYS_FILE": resolver.Config.GlobalAuthorizedKeysFile,
			"SELF_OVERRIDES_FILE":         resolver.Config.SelfOverridesFile,
			"RESOLVER_CONFIG":             resolver.ConfigPath(),
		},
		TouchFiles:       []string{resolver.ConfigName},
		CopyContextFiles: []string{core.Executable},
	})
}

func (resolver Resolver) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		core.Executable: resolver.Executable,
	}
}

func (resolver Resolver) Server(io stream.IOStream) (p wdresolve.ResolveHandler, err error) {
	p.TrustXForwardedProto = true

	fallback := &resolvers.Regexp{
		Data: map[string]string{},
	}

	// handle the default domain name!
	domainName := resolver.Config.DefaultDomain
	if domainName != "" {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domainName))] = fmt.Sprintf("https://$1.%s", domainName)
		io.Printf("registering default domain %s\n", domainName)
	}

	// handle the extra domains!
	for _, domain := range resolver.Config.SelfExtraDomains {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domain))] = fmt.Sprintf("https://$1.%s", domainName)
		io.Printf("registering legacy domain %s\n", domain)
	}

	// open the prefix file
	prefixFile := resolver.ConfigPath()
	fs, err := os.Open(prefixFile)
	io.Println("loading prefixes from ", prefixFile)
	if err != nil {
		return p, err
	}
	defer fs.Close()

	// read the prefixes
	// TODO: Do we want to load these without a file?
	prefixes, err := resolvers.ReadPrefixes(fs)
	if err != nil {
		return p, err
	}

	// and use that as the resolver!
	p.Resolver = resolvers.InOrder{
		prefixes,
		fallback,
	}

	return p, nil
}
