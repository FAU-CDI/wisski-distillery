package auth

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

// CheckScope checks if the given session has the given scope.
// If the user is denied a scope, the error will be of type AccessDeniedError.
func (auth *Auth) CheckScope(param string, scope component.Scope, r *http.Request) error {
	// get all the infos about all of the scopes
	infos := auth.scopeInfos.Get(func() []component.ScopeInfo {
		infos := make([]component.ScopeInfo, len(auth.Dependencies.ScopeProviders))
		for i, p := range auth.Dependencies.ScopeProviders {
			infos[i] = p.Scope()
		}
		return infos
	})

	// find where in teh list of parameters it is!
	index, ok := auth.scopeIndex.Get(func() map[component.Scope]int {
		m := make(map[component.Scope]int, len(infos))
		for idx, i := range infos {
			m[i.Scope] = idx
		}
		return m
	})[scope]

	if !ok {
		return ErrUnknownScope
	}

	// check that we take a parameter
	if infos[index].TakesParam && param == "" {
		return ErrParamRequired
	}
	if !infos[index].TakesParam && param != "" {
		return ErrNoParam
	}

	// call the checker and return an error
	ok, err := auth.Dependencies.ScopeProviders[index].HasScope(param, r)
	if err != nil {
		return infos[index].CheckError(err)
	}
	if ok {
		return nil
	}
	return infos[index].DeniedError()
}
