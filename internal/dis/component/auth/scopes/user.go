package scopes

import (
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
)

type UserLoggedIn struct {
	component.Base
	Dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.ScopeProvider = (*UserLoggedIn)(nil)
)

const (
	ScopeUserLoggedIn Scope = "login.user"
)

func (*UserLoggedIn) Scope() component.ScopeInfo {
	return component.ScopeInfo{
		Scope:       ScopeUserLoggedIn,
		Description: "session has an associated user",
		TakesParam:  false,
	}
}

func (iu *UserLoggedIn) HasScope(param string, r *http.Request) (bool, error) {
	user, _, err := iu.Dependencies.Auth.SessionOf(r)
	return user != nil, err
}
