package scopes

import (
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
)

type ResolverScope struct {
	component.Base
	dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.ScopeProvider = (*ResolverScope)(nil)
)

const (
	ScopeResolver Scope = "url.resolve"
)

func (*ResolverScope) Scope() component.ScopeInfo {
	return component.ScopeInfo{
		Scope:         ScopeResolver,
		Description:   "resolve a URI to a URL to display it in",
		DeniedMessage: "",
		TakesParam:    false,
	}
}

func (rs *ResolverScope) HasScope(param string, r *http.Request) (bool, error) {
	_, user, err := rs.dependencies.Auth.SessionOf(r)
	return user != nil, err
}
