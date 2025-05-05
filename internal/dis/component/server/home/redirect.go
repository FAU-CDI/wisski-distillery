//spellchecker:words home
package home

//spellchecker:words context encoding json http strings github wisski distillery internal component
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/errorsx"
)

func (home *Home) loadRedirect(context.Context) (redirect Redirect, e error) {
	if redirect.Overrides == nil {
		redirect.Overrides = make(map[string]string)
	}

	delete(redirect.Overrides, "") // make sure there is no root redirect

	redirect.Absolute = false
	redirect.Permanent = false

	// load the overrides file
	overrides, err := os.Open(component.GetStill(home).Config.Paths.OverridesJSON)
	if err != nil {
		return redirect, fmt.Errorf("failed to open overrides file: %w", err)
	}
	defer errorsx.Close(overrides, &e, "overrides file")

	// decode the overrides file
	if err := json.NewDecoder(overrides).Decode(&redirect.Overrides); err != nil {
		return redirect, fmt.Errorf("failed to parse overrides files: %w", err)
	}

	// and return!
	return redirect, nil
}

// Redirect implements a redirect server that redirects all requests.
// It implements http.Handler.
type Redirect struct {
	// Target is the target URL to redirect to.
	Target string

	// Fallback is used when target is the empty string.
	Fallback http.Handler

	// Absolute determines if the request path should be appended to the target URL when redirecting.
	// By default this path is always appended, set Absolute to true to prevent this.
	Absolute bool

	// Overrides is a map from paths to URLs that should override the default target.
	Overrides map[string]string

	// Permanent determines if the redirect responses issued should return
	// Permanent Redirect (Status Code 308) or Temporary Redirect (Status Code 307).
	Permanent bool
}

// Redirect determines the redirect URL for a specific incoming request
// If it returns the empty string, the fallback is used.
func (redirect Redirect) Redirect(r *http.Request) string {
	// if we have an override for this URL, use it immediatly
	url := strings.TrimSuffix(r.URL.Path, "/")
	if override, ok := redirect.Overrides[url]; ok {
		return override
	}

	if redirect.Target == "" {
		return ""
	}

	// if we are in absolute redirect mode, always return the absolute URL
	if redirect.Absolute {
		return redirect.Target
	}

	// return the target + the redirected URL
	dest := strings.TrimSuffix(redirect.Target, "/") + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		dest += "?" + r.URL.RawQuery
	}
	return dest
}

// ServeHTTP implements the http.Handler interface and redirects a single request to redirect.Target.
func (redirect Redirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dest := redirect.Redirect(r)
	if dest == "" {
		if redirect.Fallback == nil {
			http.NotFound(w, r)
			return
		}
		redirect.Fallback.ServeHTTP(w, r)
		return
	}

	// determine if we are temporary or permanent redirect
	status := http.StatusTemporaryRedirect
	if redirect.Permanent {
		status = http.StatusPermanentRedirect
	}

	// and do the redirect
	http.Redirect(w, r, dest, status)
}
