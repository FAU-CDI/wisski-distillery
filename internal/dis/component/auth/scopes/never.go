package scopes

import (
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
)

type Never struct {
	component.Base
	dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.ScopeProvider = (*Never)(nil)
)

const (
	ScopeNever Scope = "never"
)

func (*Never) Scope() component.ScopeInfo {
	return component.ScopeInfo{
		Scope:         ScopeNever,
		Description:   "scope that is never fullfilled",
		DeniedMessage: "no one can do this",
		TakesParam:    false,
	}
}

func (*Never) HasScope(string, *http.Request) (bool, error) {
	return false, nil
}
