package component

import (
	"html/template"
	"net/http"
)

// Menuable is a component that provides a menu
type Menuable interface {
	Component

	Menu(r *http.Request) []MenuItem
}
type MenuItem struct {
	Title  string
	Path   template.URL
	Active bool

	Priority MenuPriority // menu priority
}

// DummyMenuItem is a dummy menu item
// It should be replaced before being displayed to the user
var DummyMenuItem = MenuItem{
	Title: "* to be replaced *",
}

func MenuItemSort(a, b MenuItem) bool {
	return a.Priority < b.Priority
}

type MenuPriority int

// Menu* indicates priorities of the menu
const (
	MenuHome MenuPriority = iota
	MenuNews
	MenuResolver
	MenuUser
	MenuAdmin
	MenuAuth
)

const (
	SmallButton MenuPriority = -1
)
