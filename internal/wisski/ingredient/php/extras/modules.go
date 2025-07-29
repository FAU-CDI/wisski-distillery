//spellchecker:words extras
package extras

//spellchecker:words context embed slices strings github wisski distillery internal phpx status ingredient
import (
	"context"
	_ "embed"
	"slices"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
)

type Modules struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

var (
	_ ingredient.WissKIFetcher = (*Modules)(nil)
)

//go:embed modules.php
var modulesPHP string

// All returns the ids of all pathbuilders in consistent order.
//
// server is the server to fetch the pathbuilders from, any may be nil.
func (modules *Modules) Get(ctx context.Context, server *phpx.Server) (infos []status.DrushExtendedModuleInfo, err error) {
	err = modules.dependencies.PHP.ExecScript(ctx, server, &infos, modulesPHP, "build_extended_infos")
	if err != nil {
		return
	}

	slices.SortFunc(infos, func(left, right status.DrushExtendedModuleInfo) int {
		if left.Enabled != right.Enabled {
			if left.Enabled {
				return -1
			} else {
				return 1
			}
		}

		if types := strings.Compare(left.Type, right.Type); types != 0 {
			return types
		}

		leftHasComposer := left.HasComposer()
		rightHasComposer := right.HasComposer()
		if leftHasComposer != rightHasComposer {
			if leftHasComposer {
				return 1
			} else {
				return -1
			}
		}

		return strings.Compare(left.Name, right.Name)
	})
	return
}

func (modules *Modules) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.Modules, _ = modules.Get(flags.Context, flags.Server)
	return
}
