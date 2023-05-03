package component

import (
	"fmt"
	"net/http"
)

// Scope represents a single permit by a session to perform some action.
// Scopes consist of two parts: A general name and a specific object.
// Scopes are checked by ScopeCheckers.
type Scope string

type ScopeInfo struct {
	Scope

	// Description is a human readable description of the scope
	Description string

	// error returned to a user when the permission is denied.
	// defaults to "missing scope {{ name }}"
	DeniedMessage string

	// TakesParam indicates if the scope accepts a parameter
	TakesParam bool
}

type CheckError struct {
	Scope Scope
	Err   error
}

func (ce CheckError) Unwrap() error { return ce.Err }
func (ce CheckError) Error() string {
	return fmt.Sprintf("unable to check scope %q: %s", string(ce.Scope), ce.Err)
}

type AccessDeniedError string

func (aed AccessDeniedError) Error() string { return string(aed) }

// DeniedError returns an AccessDeniedError that indivates the access is denied.
func (scope ScopeInfo) DeniedError() error {
	if scope.DeniedMessage == "" {
		return AccessDeniedError(fmt.Sprintf("missing scope %q", string(scope.Scope)))
	}
	return AccessDeniedError(scope.DeniedMessage)
}

// CheckError returns a CheckError with the given underlying error.
func (scope ScopeInfo) CheckError(err error) error {
	return CheckError{Scope: scope.Scope, Err: err}
}

const (
	ScopeUserLoggedIn  Scope = "login.user"
	ScopeAdminLoggedIn Scope = "login.admin"
)

// ScopeProvider is a component that can check a specific scope
type ScopeProvider interface {
	Component

	// Scopes returns information about the scope
	Scope() ScopeInfo

	// Check checks if the given session has access to the given scope
	HasScope(param string, r *http.Request) (bool, error)
}
