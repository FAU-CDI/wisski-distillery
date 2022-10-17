package wisski

import (
	"fmt"
	"net/url"
)

// Domain returns the full domain name of this WissKI
func (wisski WissKI) Domain() string {
	return fmt.Sprintf("%s.%s", wisski.Slug, wisski.Core.Config.DefaultDomain)
}

// URL returns the public URL of this instance
func (wisski WissKI) URL() *url.URL {
	// setup domain and path
	url := &url.URL{
		Host: wisski.Domain(),
		Path: "/",
	}

	// use http or https scheme depending on if the distillery has it enabled
	if wisski.Core.Config.HTTPSEnabled() {
		url.Scheme = "https"
	} else {
		url.Scheme = "http"
	}

	return url
}
