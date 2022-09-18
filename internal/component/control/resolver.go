package control

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/FAU-CDI/wdresolve"
	"github.com/FAU-CDI/wdresolve/resolvers"
	"github.com/tkw1536/goprogram/stream"
)

func (control Control) ResolverConfigPath() string {
	return filepath.Join(control.Dir, control.ResolverFile)
}

func (control Control) resolver(io stream.IOStream) (p wdresolve.ResolveHandler, err error) {
	p.TrustXForwardedProto = true

	fallback := &resolvers.Regexp{
		Data: map[string]string{},
	}

	// handle the default domain name!
	domainName := control.Config.DefaultDomain
	if domainName != "" {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domainName))] = fmt.Sprintf("https://$1.%s", domainName)
		io.Printf("registering default domain %s\n", domainName)
	}

	// handle the extra domains!
	for _, domain := range control.Config.SelfExtraDomains {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domain))] = fmt.Sprintf("https://$1.%s", domainName)
		io.Printf("registering legacy domain %s\n", domain)
	}

	// open the prefix file
	prefixFile := control.ResolverConfigPath()
	fs, err := control.Environment.Open(prefixFile)
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
