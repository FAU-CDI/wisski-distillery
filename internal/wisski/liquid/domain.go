package liquid

import (
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
)

// Domain returns the full domain name of this WissKI
func (liquid *Liquid) Domain() string {
	return liquid.Config.HTTP.HostFromSlug(liquid.Slug)
}

func (liquid *Liquid) Hostname() string {
	return liquid.Domain() + ".wisski"
}

// HostRule returns a host rule for this wisski
func (liquid *Liquid) HostRule() string {
	return config.MakeHostRule(liquid.Domain())
}

// URL returns the public URL of this instance
func (liquid *Liquid) URL() *url.URL {
	// setup domain and path
	url := &url.URL{
		Host: liquid.Domain(),
		Path: "/",
	}

	// use http or https scheme depending on if the distillery has it enabled
	if liquid.Malt.Config.HTTP.HTTPSEnabled() {
		url.Scheme = "https"
	} else {
		url.Scheme = "http"
	}

	return url
}
