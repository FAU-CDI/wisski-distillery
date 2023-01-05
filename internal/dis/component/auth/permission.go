package auth

import (
	"errors"
	"net/http"
)

// Permission represents a permission granted to a user.
//
// The nil permission represents any authenticated user.
type Permission func(user *AuthUser, r *http.Request) (ok Grant, err error)

// Grant represents an object that either grants or denies access for a certain permission
type Grant interface {
	isGranted()

	// Granted returns a boolean indicating if permission to the resource in question
	// has been granted
	Granted() bool

	// Denied returns a string containing an error message to display to the user when permission is denied.
	// When Granted() returns true, the behaviour is undefined.
	Denied() string
}

// Bool2Grant returns a new grant that returns granted for the given boolean, and message as the denied message.
func Bool2Grant(granted bool, message string) Grant {
	if granted {
		return grantAllow{}
	}
	return grantDeny(message)
}

type grantAllow struct{}

func (grantAllow) isGranted()     {}
func (grantAllow) Granted() bool  { return true }
func (grantAllow) Denied() string { return "" }

type grantDeny string

func (grantDeny) isGranted()      {}
func (g grantDeny) Granted() bool { return false }
func (g grantDeny) Denied() string {
	if g == "" {
		return "Forbidden"
	}
	return string(g)
}

// AllPermissions returns a new permission that checks if all the given permissions are set
func AllPermissions(clauses ...Permission) Permission {
	return func(user *AuthUser, r *http.Request) (ok Grant, err error) {
		for _, clause := range clauses {
			perm, err := clause.Permit(user, r)
			if err != nil {
				return perm, err
			}
			if !perm.Granted() {
				return perm, nil
			}
		}

		// everything was fine
		return grantAllow{}, nil
	}
}

var errPermissionPanic = errors.New("permission: panic()")

// Permit checks if the given user has this permission.
func (perm Permission) Permit(user *AuthUser, r *http.Request) (ok Grant, err error) {
	// if there is no permission, then we just check if there is some user
	if perm == nil {
		return Bool2Grant(user != nil, ""), nil
	}

	// recover any panic()ed permission call
	// to prevent the handler from panic()ing
	defer func() {
		if p := recover(); p != nil {
			ok = Bool2Grant(false, "unknown error")
			err = errPermissionPanic
		}
	}()

	return perm(user, r)
}
