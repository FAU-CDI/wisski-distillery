//spellchecker:words component
package component

//spellchecker:words html template http sync atomic
import (
	"html/template"
	"net/http"
	"sync/atomic"
)

// Menuable is a component that provides a menu.
type Menuable interface {
	Component

	Menu(r *http.Request) []MenuItem
}

// MenuItem represents an item inside the menu.
type MenuItem struct {
	Title  string
	Path   template.URL
	Active bool // Active, only used for tabs and crumbs
	Sticky bool // Sticky, and do not collapse when collapsing the menu (ignored for tabs and crumbs)

	Priority MenuPriority

	replaceID uint64 // internal id used to replace an item
}

var dummyCounter uint64

// DummyMenuItem creates a new Dummy Menu Item to be replaced.
func DummyMenuItem() MenuItem {
	return MenuItem{
		replaceID: atomic.AddUint64(&dummyCounter, 1),
	}
}

// ReplaceWith replaces this MenuItem with a different MenuItem.
// This method returns true if an appropriate DummyMenuItem exists.
func (mi MenuItem) ReplaceWith(new MenuItem, items []MenuItem) bool {
	if mi.replaceID == 0 {
		// never replace non-dummy items
		return false
	}
	for i, item := range items {
		if mi.replaceID == item.replaceID {
			items[i] = new
			return true
		}
	}

	return false
}

func MenuItemSort(a, b MenuItem) int {
	return int(a.Priority) - int(b.Priority)
}

type MenuPriority int

// Menu* indicates priorities of the menu.
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
