package custom

import (
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/mux"
	"golang.org/x/exp/slices"
)

// getMenuItems gets a fresh copy of the cached slice of menu items
func (custom *Custom) Menu(r *http.Request) []component.MenuItem {
	return custom.menu.Get(func() []component.MenuItem {
		items := make([]component.MenuItem, 0, len(custom.Dependencies.Routeables))
		for _, route := range custom.Dependencies.Routeables {
			routes := route.Routes()
			if routes.MenuTitle == "" {
				continue
			}
			items = append(items, component.MenuItem{
				Title:    routes.MenuTitle,
				Priority: routes.MenuPriority,
				Path:     template.URL(routes.Prefix),
			})
		}
		slices.SortFunc(items, component.MenuItemSort)
		return items
	})
}

func (custom *Custom) BuildMenu(r *http.Request) []component.MenuItem {
	// NOTE(twiesing): Don't name this method "Menu", as it will cause
	// a stack overflow.
	path := mux.NormalizePath(r.URL.Path)

	// get the static menu items, and then return all the regular ones
	var items []component.MenuItem
	for _, m := range custom.Dependencies.Menuable {
		items = append(items, m.Menu(r)...)
	}
	for i, item := range items {
		items[i].Active = string(item.Path) == path
	}
	slices.SortFunc(items, component.MenuItemSort)
	return items
}
