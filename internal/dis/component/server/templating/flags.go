package templating

import (
	"html/template"
	"net/http"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/tkw1536/goprogram/lib/reflectx"
	"golang.org/x/exp/slices"
)

// Flags represent handle-updatable options for the base template
type Flags struct {
	Title         string // Title of the menu
	assets.Assets        // assets are the assets included in the template

	Crumbs  []component.MenuItem // crumbs are the breadcrumbs leading to a specific action
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
// Any context may e
type RuntimeFlags struct {
	Flags

	RequestURI  string               // request uri of the current page
	Menu        []component.MenuItem // menu at the top of the page
	GeneratedAt time.Time            // time the underlying data returned
	CSRF        template.HTML        // csrf data (if any)
}

var runtimeFlagsName = reflectx.TypeOf[RuntimeFlags]().Name()

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
		flags.Crumbs = crumbs
		return flags
	}
}

// Actions sets the actions
func Actions(actions ...component.MenuItem) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		flags.Actions = actions
		return flags
	}
}

// ReplaceAction replaces a specific action
func ReplaceAction(index int, action component.MenuItem) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		flags.Actions[index] = action
		return flags
	}
}

// ReplaceCrumb replaces a specific crum
func ReplaceCrumb(index int, action component.MenuItem) FlagFunc {
	return func(flags Flags, r *http.Request) Flags {
		flags.Crumbs[index] = action
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
