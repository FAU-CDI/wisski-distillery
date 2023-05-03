package scopes

import (
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
)

type AdminLoggedIn struct {
	component.Base
	Dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.ScopeProvider = (*UserLoggedIn)(nil)
)

func (*AdminLoggedIn) Scope() component.ScopeInfo {
	return component.ScopeInfo{
		Scope:         component.ScopeAdminLoggedIn,
		Description:   "session has a signed in admin",
		DeniedMessage: "user must be signed into an admin account with TOTP enabled",
		TakesParam:    false,
	}
}

func (al *AdminLoggedIn) HasScope(param string, r *http.Request) (bool, error) {
	user, err := al.Dependencies.Auth.UserOf(r)
	return user != nil && user.IsAdmin() && user.IsTOTPEnabled(), err
}
