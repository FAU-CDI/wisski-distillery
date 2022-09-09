package env

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/FAU-CDI/wdresolve"
	"github.com/FAU-CDI/wdresolve/resolvers"
	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
	"github.com/tkw1536/goprogram/stream"
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
	return resolver.dis.makeComponentStack(resolver, stack.Installable{
		EnvFileContext: map[string]string{
			"VIRTUAL_HOST":      resolver.dis.DefaultVirtualHost(),
			"LETSENCRYPT_HOST":  resolver.dis.DefaultLetsencryptHost(),
			"LETSENCRYPT_EMAIL": resolver.dis.Config.CertbotEmail,

			"CONFIG_PATH": resolver.dis.Config.ConfigPath,
			"DEPLOY_ROOT": resolver.dis.Config.DeployRoot,

			"GLOBAL_AUTHORIZED_KEYS_FILE": resolver.dis.Config.GlobalAuthorizedKeysFile,
			"SELF_OVERRIDES_FILE":         resolver.dis.Config.SelfOverridesFile,
			"RESOLVER_CONFIG":             resolver.ConfigPath(),
		},
		CopyContextFiles: []string{core.Executable},
	})
}

func (resolver ResolverComponent) Context(parent stack.InstallationContext) stack.InstallationContext {
	return stack.InstallationContext{
		core.Executable: resolver.dis.CurrentExecutable(),
	}
}

func (resolver ResolverComponent) Server(io stream.IOStream) (p wdresolve.ResolveHandler, err error) {
	p.TrustXForwardedProto = true

	fallback := &resolvers.Regexp{
		Data: map[string]string{},
	}

	// handle the default domain name!
	domainName := resolver.dis.Config.DefaultDomain
	if domainName != "" {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domainName))] = fmt.Sprintf("https://$1.%s", domainName)
		io.Printf("registering default domain %s\n", domainName)
	}

	// handle the extra domains!
	for _, domain := range resolver.dis.Config.SelfExtraDomains {
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

func (resolver ResolverComponent) Path() string {
	return resolver.dis.getComponentPath(resolver)
}

func (resolver ResolverComponent) ConfigPath() string {
	return filepath.Join(resolver.Path(), resolver.ConfigName)
}
