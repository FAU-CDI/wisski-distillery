//spellchecker:words extras
package extras

//spellchecker:words context embed slices github wisski distillery internal phpx status ingredient
import (
	"context"
	_ "embed"
	"slices"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
)

type Modules struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

//go:embed modules.php
var modulesPHP string

type DrushExtendedModuleInfo struct {
	DrushModuleInfo
	Composer *ComposerModuleInfo `json:"composer"`
}

func (demi DrushExtendedModuleInfo) HasComposer() bool {
	return demi.Composer != nil
}

type DrushModuleInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`

	Type    string `json:"type"`
	Path    string `json:"path"`
	Enabled bool   `json:"enabled"`
	Version string `json:"version"`
}

type ComposerModuleInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

// All returns the ids of all pathbuilders in consistent order.
//
// server is the server to fetch the pathbuilders from, any may be nil.
func (modules *Modules) Get(ctx context.Context, server *phpx.Server) (infos []DrushExtendedModuleInfo, err error) {
	err = modules.dependencies.PHP.ExecScript(ctx, server, &infos, modulesPHP, "build_extended_infos")
	if err != nil {
		return
	}

	slices.SortFunc(infos, func(left, right DrushExtendedModuleInfo) int {
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
