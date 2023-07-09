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

func (wisski *WissKI) allIngredients() []initFunc {
	return []initFunc{
		// core bits
		auto[*locker.Locker],
		manual(func(m *mstore.MStore) {
			m.Storage = wisski.Malt.Meta.Storage(wisski.Slug)
		}),

		// php
		auto[*php.PHP],
		auto[*extras.Prefixes],
		auto[*extras.Settings],
		auto[*extras.Pathbuilder],
		auto[*extras.Stats],
		auto[*extras.Blocks],
		auto[*extras.Requirements],
		auto[*extras.Adapters],
		auto[*users.Users],
		auto[*users.UserPolicy],

		// info
		manual(func(info *info.Info) {
			info.Analytics = &wisski.lifetime.Analytics
		}),
		auto[*barrel.LastRebuildFetcher],
		auto[*barrel.RunningFetcher],
		auto[*drush.LastUpdateFetcher],
		auto[*drush.LastCronFetcher],
		auto[*info.SnapshotsFetcher],

		// stacks
		auto[*barrel.Barrel],
		auto[*bookkeeping.Bookkeeping],
		auto[*manager.Manager],
		auto[*system.SystemManager],
		auto[*composer.Composer],
		auto[*drush.Drush],

		auto[*reserve.Reserve],

		auto[*ssh.SSH],
	}
}
