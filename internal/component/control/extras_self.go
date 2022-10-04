package control

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/tkw1536/goprogram/stream"
)

// SelfHandler implements serving the '/' route
type SelfHandler struct {
	component.ComponentBase

	Instances *instances.Instances
}

func (SelfHandler) Name() string { return "control-self" }

func (*SelfHandler) Routes() []string { return []string{"/"} }

func (sh *SelfHandler) Handler(route string, io stream.IOStream) (http.Handler, error) {
	// create a redirect
	var redirect Redirect
	var err error

	// open the overrides file
	overrides, err := sh.Environment.Open(sh.Config.SelfOverridesFile)
	io.Printf("loading overrides from %q\n", sh.Config.SelfOverridesFile)
	if err != nil {
		return redirect, err
	}
	defer overrides.Close()

	// decode the overrides file
	if err := json.NewDecoder(overrides).Decode(&redirect.Overrides); err != nil {
		return nil, err
	}

	if redirect.Overrides == nil {
		redirect.Overrides = make(map[string]string)
	}
	redirect.Overrides[""] = sh.Config.SelfRedirect.String()

	// create a redirect server
	redirect.Fallback, err = sh.selfFallback()
	if err != nil {
		return nil, err
	}
	redirect.Absolute = false
	redirect.Permanent = false

	// and return!
	return redirect, nil
}

func (sh *SelfHandler) selfFallback() (http.Handler, error) {
	return http.HandlerFunc(sh.serveFallback), nil
}

var notFoundText = []byte("not found")

func (sh *SelfHandler) serveFallback(w http.ResponseWriter, r *http.Request) {

	slug := sh.Config.SlugFromHost(r.Host)
	if slug == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write(notFoundText)
		return
	}

	if ok, _ := sh.Instances.Has(slug); !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "WissKI %q not found\n", slug)
		return
	}

	w.WriteHeader(http.StatusBadGateway)
	fmt.Fprintf(w, "WissKI %q is currently offline\n", slug)

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
