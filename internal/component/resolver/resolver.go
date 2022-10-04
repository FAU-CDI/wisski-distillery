package resolver

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/FAU-CDI/wdresolve"
	"github.com/FAU-CDI/wdresolve/resolvers"
	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/control"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/tkw1536/goprogram/stream"
)

type Resolver struct {
	component.ComponentBase

	Control      *control.Control
	ResolverFile string

	handler lazy.Lazy[wdresolve.ResolveHandler]
}

func (Resolver) Name() string { return "resolver" }

func (resolver *Resolver) Routes() []string { return []string{"/go/", "/wisski/get/"} }

func (resolver *Resolver) Handler(route string, io stream.IOStream) (http.Handler, error) {
	var err error
	return resolver.handler.Get(func() (p wdresolve.ResolveHandler) {
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

		configPath := resolver.ConfigPath()
		{
			// load the prefix path!
			var fs fs.File
			fs, err = resolver.Environment.Open(configPath)
			io.Println("loading prefixes from ", configPath)
			if err != nil {
				return
			}
			defer fs.Close()

			// read the file
			var prefixes resolvers.Prefix
			prefixes, err = resolvers.ReadPrefixes(fs)
			if err != nil {
				return
			}

			// and use that as the resolver!
			p.Resolver = resolvers.InOrder{
				prefixes,
				fallback,
			}

			return p
		}
	}), err
}

func (resolver *Resolver) ConfigPath() string {
	return filepath.Join(resolver.Control.Path(), resolver.ResolverFile)
}
