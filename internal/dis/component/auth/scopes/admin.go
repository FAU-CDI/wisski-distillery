//spellchecker:words scopes
package scopes

//spellchecker:words http github wisski distillery internal component auth
import (
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
)

type AdminLoggedIn struct {
	component.Base
	dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.ScopeProvider = (*UserLoggedIn)(nil)
)

const (
	ScopeUserAdmin Scope = "user.admin"
)

func (*AdminLoggedIn) Scope() component.ScopeInfo {
	return component.ScopeInfo{
		Scope:         ScopeUserAdmin,
		Description:   "session must have a valid admin",
		DeniedMessage: "user must have an admin account with TOTP enabled",
		TakesParam:    false,
	}
}

func (al *AdminLoggedIn) HasScope(param string, r *http.Request) (bool, error) {
	_, user, err := al.dependencies.Auth.SessionOf(r)
	return user != nil && user.IsAdmin() && user.IsTOTPEnabled(), err
}
