package liquid

import (
	"fmt"
	"net/url"
)

// Domain returns the full domain name of this WissKI
func (liquid *Liquid) Domain() string {
	return fmt.Sprintf("%s.%s", liquid.Slug, liquid.Malt.Config.DefaultDomain)
}

// URL returns the public URL of this instance
func (liquid *Liquid) URL() *url.URL {
	// setup domain and path
	url := &url.URL{
		Host: liquid.Domain(),
		Path: "/",
	}

	// use http or https scheme depending on if the distillery has it enabled
	if liquid.Malt.Config.HTTPSEnabled() {
		url.Scheme = "https"
	} else {
		url.Scheme = "http"
	}

	return url
}