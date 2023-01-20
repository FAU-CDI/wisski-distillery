package templating

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

// Templating implements templating customization
type Templating struct {
	component.Base
	Dependencies struct {
		Routeables []component.Routeable
		Menuable   []component.Menuable
	}
	menu lazy.Lazy[[]component.MenuItem]
}

var (
	_ component.Backupable = (*Templating)(nil)
	_ component.Menuable   = (*Templating)(nil)
)
