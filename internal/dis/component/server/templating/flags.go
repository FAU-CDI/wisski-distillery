package templating

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"golang.org/x/exp/slices"
)

// Flags represent handle-updatable options for the base template
type Flags struct {
	Title         string // Title of the menu
	assets.Assets        // assets are the assets included in the template

	Crumbs  []component.MenuItem // crumbs are the breadcrumbs leading to a specific action
	Tabs    []component.MenuItem // tabs are shown above actions, and act as a menu
	Actions []component.MenuItem // actions are the actions available to a specific thingy
}

// Apply applies a set of functions to this flags
func (flags Flags) Apply(r *http.Request, funcs ...FlagFunc) Flags {
	for _, f := range funcs {
		flags = f(flags, r)
	}
	return flags
}

// RuntimeFlags are passed to the template at runtime.
type RuntimeFlags struct {
	Flags

	RequestURI  string               // request uri of the current page
	Menu        []component.MenuItem // menu at the top of the page
	StartedAt   time.Time            // time the request started to generate
	GeneratedAt time.Time            // time the underlying data returned
	CSRF        template.HTML        // csrf data (if any)
}

// Returns how long this request took to render
func (rf RuntimeFlags) Took() time.Duration {
	return time.Since(rf.StartedAt)
}
func (rf RuntimeFlags) TookHTML() template.HTML {
	took := rf.Took()

	return template.HTML(fmt.Sprintf("<time datetime=\"P%.3f\">%s</time>", took.Seconds(), took))
}

var runtimeFlagsName = reflect.TypeFor[RuntimeFlags]().Name()

// Clone clones this flags
func (flags Flags) Clone() Flags {
	flags.Crumbs = slices.Clone(flags.Crumbs)
	flags.Actions = slices.Clone(flags.Actions)
	return flags
}

// FlagFunc updates a flags based on a request.
// FlagFunc may not be nil.
type FlagFunc func(flags Flags, r *http.Request) Flags

// Assets sets the given assets for the given flags
func Assets(Assets assets.Assets) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		flags.Assets = Assets
		return flags
	}
}

// Crumbs sets the crumbs
func Crumbs(crumbs ...component.MenuItem) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		flags.Crumbs = slices.Clone(crumbs)
		return flags
	}
}

// Actions sets the actions
func Actions(actions ...component.MenuItem) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		flags.Actions = slices.Clone(actions)
		return flags
	}
}

// Tabs sets the tabs
func Tabs(actions ...component.MenuItem) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		flags.Tabs = slices.Clone(actions)
		return flags
	}
}

// ReplaceCrumb replaces a specific crum
func ReplaceCrumb(old component.MenuItem, crumb component.MenuItem) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		if !old.ReplaceWith(crumb, flags.Crumbs) {
			wdlog.Of(r.Context()).Warn(
				"did not replace crumb item",
				"crum", fmt.Sprint(crumb),
				"crumbs", fmt.Sprint(flags.Crumbs),
			)
		}
		return flags
	}
}

// ReplaceAction replaces a specific action
func ReplaceAction(old component.MenuItem, action component.MenuItem) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		if !old.ReplaceWith(action, flags.Actions) {
			wdlog.Of(r.Context()).Warn(
				"did not replace action item",
				"action", fmt.Sprint(action),
				"actions", fmt.Sprint(flags.Actions),
			)
		}
		return flags
	}
}

// ReplaceTab replaces a specific tab
func ReplaceTab(old component.MenuItem, tab component.MenuItem) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		if !old.ReplaceWith(tab, flags.Tabs) {
			wdlog.Of(r.Context()).Warn(
				"did not replace menu item",
				"tab", fmt.Sprint(tab),
				"tabs", fmt.Sprint(flags.Tabs),
			)
		}
		return flags
	}
}

// Title sets the title of this template
func Title(title string) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		flags.Title = title
		return flags
	}
}
