//spellchecker:words scopes
package scopes

//spellchecker:words http github wisski distillery internal component auth
import (
	"fmt"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
)

type ListInstancesScope struct {
	component.Base
	dependencies struct {
		Auth *auth.Auth
	}
}

var (
	_ component.ScopeProvider = (*ListInstancesScope)(nil)
)

const (
	ScopeInstanceDirectory Scope = "instances.directory"
)

func (*ListInstancesScope) Scope() component.ScopeInfo {
	return component.ScopeInfo{
		Scope:         ScopeInstanceDirectory,
		Description:   "get a public directory of instances",
		DeniedMessage: "",
		TakesParam:    false,
	}
}

func (lis *ListInstancesScope) HasScope(param string, r *http.Request) (bool, error) {
	_, user, err := lis.dependencies.Auth.SessionOf(r)
	if err != nil {
		return false, fmt.Errorf("failed to get session: %w", err)
	}
	return user != nil, nil
}
