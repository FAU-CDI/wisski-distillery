//spellchecker:words auth
package auth

//spellchecker:words errors http github wisski distillery internal component
import (
	"errors"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

var (
	ErrUnknownScope  = errors.New("unknown scope")
	ErrParamRequired = errors.New("scope requires parameter")
	ErrNoParam       = errors.New("scope does not take parameter")
)

// Scopes returns a map of all available scopes.
func (auth *Auth) Scopes() map[component.Scope]component.ScopeInfo {
	scopes := auth.getScopeMap()
	mp := make(map[component.Scope]component.ScopeInfo, len(scopes))
	for scope, entry := range scopes {
		mp[scope] = entry.Info
	}
	return mp
}

// getScopeMap return a (cached version of) all scopes.
func (auth *Auth) getScopeMap() map[component.Scope]scopeMapEntry {
	return auth.scopeMap.Get(func() map[component.Scope]scopeMapEntry {
		mp := make(map[component.Scope]scopeMapEntry, len(auth.dependencies.ScopeProviders))
		for _, p := range auth.dependencies.ScopeProviders {
			info := p.Scope()
			mp[info.Scope] = scopeMapEntry{
				Provider: p,
				Info:     info,
			}
		}
		return mp
	})
}

// CheckScope checks if the given request is associated with the given request.
// A request can be one of two types:
// - A signed in user with an implicitly associated set of scopes
// - A session authorized with a token only
// If the request is denied a scope, the error will be of type AccessDeniedError.
func (auth *Auth) CheckScope(param string, scope component.Scope, r *http.Request) error {
	// the empty scope is always permitted implicitly
	if scope == "" {
		return nil
	}

	entry, ok := auth.getScopeMap()[scope]
	if !ok {
		return ErrUnknownScope
	}

	// check that we take a parameter
	if entry.Info.TakesParam && param == "" {
		return ErrParamRequired
	}
	if !entry.Info.TakesParam && param != "" {
		return ErrNoParam
	}

	// call the checker and return an error
	ok, err := entry.Provider.HasScope(param, r)
	if err != nil {
		return entry.Info.CheckError(err)
	}
	if ok {
		return nil
	}
	return entry.Info.DeniedError()
}
