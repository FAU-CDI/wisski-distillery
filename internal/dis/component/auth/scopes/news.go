//spellchecker:words scopes
package scopes

//spellchecker:words http github wisski distillery internal component auth
import (
	"fmt"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
)

type ListNewsScope struct {
	component.Base
	dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.ScopeProvider = (*ListNewsScope)(nil)
)

const (
	ScopeListNews Scope = "news.list"
)

func (*ListNewsScope) Scope() component.ScopeInfo {
	return component.ScopeInfo{
		Scope:         ScopeListNews,
		Description:   "list news items",
		DeniedMessage: "",
		TakesParam:    false,
	}
}

func (lns *ListNewsScope) HasScope(param string, r *http.Request) (bool, error) {
	_, user, err := lns.dependencies.Auth.SessionOf(r)
	if err != nil {
		return false, fmt.Errorf("failed to get session: %w", err)
	}
	return user != nil, nil
}
