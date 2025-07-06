//spellchecker:words templating
package templating

//spellchecker:words html template http slices github wisski distillery internal component pkglib httpx
import (
	"html/template"
	"net/http"
	"slices"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"go.tkw01536.de/pkglib/httpx/mux"
)

// buildMenu builds the manu for this request for all known components in this distillery.
//
// NOTE(twiesing): Don't name this method "Menu", as it will cause a stack overflow.
func (tpl *Templating) buildMenu(r *http.Request) []component.MenuItem {
	path := mux.NormalizePath(r.URL.Path)

	// get the static menu items, and then return all the regular ones
	var items []component.MenuItem
	for _, m := range tpl.dependencies.Menuable {
		items = append(items, m.Menu(r)...)
	}
	for i, item := range items {
		items[i].Active = string(item.Path) == path
	}
	slices.SortFunc(items, component.MenuItemSort)
	return items
}

// Menu returns a list of menu items provided by routeables.
func (tpl *Templating) Menu(r *http.Request) []component.MenuItem {
	return tpl.menu.Get(func() []component.MenuItem {
		items := make([]component.MenuItem, 0, len(tpl.dependencies.Routeables))
		for _, route := range tpl.dependencies.Routeables {
			routes := route.Routes()
			if routes.MenuTitle == "" {
				continue
			}
			items = append(items, component.MenuItem{
				Title:    routes.MenuTitle,
				Priority: routes.MenuPriority,
				Sticky:   routes.MenuSticky,
				Path:     template.URL(routes.Prefix), // #nosec G203 -- prefix assumed to be safe
			})
		}
		slices.SortFunc(items, component.MenuItemSort)
		return items
	})
}
