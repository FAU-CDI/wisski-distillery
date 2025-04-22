//spellchecker:words component
package component

//spellchecker:words encoding json http github wisski distillery internal models
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
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

// ScopeProvider is a component that can check a specific scope.
type ScopeProvider interface {
	Component

	// Scopes returns information about the scope
	Scope() ScopeInfo

	// Check checks if the given session has access to the given scope.
	HasScope(param string, r *http.Request) (bool, error)
}

// SessionInfo provides information about the current session.
type SessionInfo struct {
	// User is the current user associated with the session.
	User *models.User

	// Token indicates if the user was authenticated with a Token
	Token bool
}

func (si SessionInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		User  string `json:"user"`
		Token bool   `json:"token"`
	}{User: si.Username(), Token: si.Token})
}

// Username reports the username associated with this session.
func (si SessionInfo) Username() string {
	if si.User == nil {
		return ""
	}
	return si.User.User
}

// Anonymous reports if this Session is associated with a user account.
func (si SessionInfo) Anonymous() bool {
	return si.Username() != ""
}
