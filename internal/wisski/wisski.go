// Package wisski provides WissKI
package wisski

import (
	"sync"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/composer"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/drush"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/manager"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/ssh"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/system"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/info"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/users"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/reserve"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/liquid"
	"github.com/tkw1536/pkglib/lifetime"
)

// WissKI represents a single WissKI Instance.
// A WissKI may not be copied
type WissKI struct {
	liquid.Liquid

	lifetimeInit sync.Once
	lifetime     lifetime.Lifetime[ingredient.Ingredient, *liquid.Liquid]
}

//
//  INIT & EXPORT
//

func (wisski *WissKI) init() {
	wisski.lifetimeInit.Do(func() {
		wisski.lifetime.Init = ingredient.Init
		wisski.lifetime.Register = wisski.allIngredients
	})
}

func export[I ingredient.Ingredient](wisski *WissKI) I {
	wisski.init()
	return lifetime.Export[I](&wisski.lifetime, &wisski.Liquid)
}

//lint:ignore U1000 for future use
func exportAll[I ingredient.Ingredient](wisski *WissKI) []I {
	wisski.init()
	return lifetime.ExportSlice[I](&wisski.lifetime, &wisski.Liquid)
}

//
// PUBLIC INGREDIENT GETTERS
//

func (wisski *WissKI) Locker() *locker.Locker {
	return export[*locker.Locker](wisski)
}

func (wisski *WissKI) Reserve() *reserve.Reserve {
	return export[*reserve.Reserve](wisski)
}

func (wisski *WissKI) Barrel() *barrel.Barrel {
	return export[*barrel.Barrel](wisski)
}

func (wisski *WissKI) Manager() *manager.Manager {
	return export[*manager.Manager](wisski)
}

func (wisski *WissKI) SystemManager() *system.SystemManager {
	return export[*system.SystemManager](wisski)
}

func (wisski *WissKI) PHP() *php.PHP {
	return export[*php.PHP](wisski)
}

func (wisski *WissKI) Bookkeeping() *bookkeeping.Bookkeeping {
	return export[*bookkeeping.Bookkeeping](wisski)
}

func (wisski *WissKI) Drush() *drush.Drush {
	return export[*drush.Drush](wisski)
}

func (wisski *WissKI) Composer() *composer.Composer {
	return export[*composer.Composer](wisski)
}

func (wisski *WissKI) Users() *users.Users {
	return export[*users.Users](wisski)
}

func (wisski *WissKI) Prefixes() *extras.Prefixes {
	return export[*extras.Prefixes](wisski)
}

func (wisski *WissKI) Settings() *extras.Settings {
	return export[*extras.Settings](wisski)
}

func (wisski *WissKI) Pathbuilder() *extras.Pathbuilder {
	return export[*extras.Pathbuilder](wisski)
}

func (wisski *WissKI) Info() *info.Info {
	return export[*info.Info](wisski)
}

func (wisski *WissKI) SSH() *ssh.SSH {
	return export[*ssh.SSH](wisski)
}

func (wisski *WissKI) Blocks() *extras.Blocks {
	return export[*extras.Blocks](wisski)
}

//
// All components
// THESE SHOULD NEVER BE CALLED DIRECTLY
//

func (wisski *WissKI) allIngredients(context *lifetime.RegisterContext[ingredient.Ingredient, *liquid.Liquid]) {
	// core bits
	lifetime.Place[*locker.Locker](context)
	lifetime.Register(context, func(m *mstore.MStore, _ *liquid.Liquid) {
		m.Storage = wisski.Malt.Meta.Storage(wisski.Slug)
	})

	// php
	lifetime.Place[*php.PHP](context)
	lifetime.Place[*extras.Prefixes](context)
	lifetime.Place[*extras.Settings](context)
	lifetime.Place[*extras.Pathbuilder](context)
	lifetime.Place[*extras.Stats](context)
	lifetime.Place[*extras.Blocks](context)
	lifetime.Place[*extras.Requirements](context)
	lifetime.Place[*extras.Adapters](context)
	lifetime.Place[*extras.Theme](context)
	lifetime.Place[*extras.Version](context)
	lifetime.Place[*users.Users](context)
	lifetime.Place[*users.UserPolicy](context)

	// info
	lifetime.Place[*info.Info](context)
	lifetime.Place[*barrel.LastRebuildFetcher](context)
	lifetime.Place[*barrel.RunningFetcher](context)
	lifetime.Place[*composer.LastUpdateFetcher](context)
	lifetime.Place[*drush.LastCronFetcher](context)
	lifetime.Place[*info.SnapshotsFetcher](context)

	// stacks
	lifetime.Place[*barrel.Barrel](context)
	lifetime.Place[*bookkeeping.Bookkeeping](context)
	lifetime.Place[*manager.Manager](context)
	lifetime.Place[*system.SystemManager](context)
	lifetime.Place[*composer.Composer](context)
	lifetime.Place[*drush.Drush](context)

	lifetime.Place[*reserve.Reserve](context)

	lifetime.Place[*ssh.SSH](context)
}
