//spellchecker:words templating
package templating

//spellchecker:words github wisski distillery internal component pkglib lazy
import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"go.tkw01536.de/pkglib/lazy"
)

// Templating implements templating customization.
type Templating struct {
	component.Base
	dependencies struct {
		Routeables []component.Routeable
		Menuable   []component.Menuable
	}
	menu lazy.Lazy[[]component.MenuItem]
}

var (
	_ component.Backupable = (*Templating)(nil)
	_ component.Menuable   = (*Templating)(nil)
)
