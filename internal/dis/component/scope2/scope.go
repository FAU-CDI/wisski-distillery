// Package scope2 implements version 2 of scoping.
package scope2

import "strings"

// Scope represents a permission of some session, e.g. a browser session or a user session.
// The zero value for scope it never fullfilled.
//
// Scopes each have a Kind, and an optionally associated instance.
// See the [Split] method.
//
// Scopes are validated at runtime using the auth component.
type Scope string

// Split splits this scope into its' kind and instance.
// The kind is guaranteed to be valid.
// If the specific kind does not take a parameter, or the parameter is empty, returns the empty instance.
func (scope Scope) Split() (kind Kind, instance string) {
	knd, instance, ok := strings.Cut(string(scope), ":")

	// get info about the kind
	info, valid := kindInfos[Kind(knd)]
	if !valid {
		return KindNever, ""
	}

	// no instance provided or required => return only the kind
	if !ok || instance == "" || !info.NeedsInstance {
		return Kind(knd), ""
	}

	// it was valid => we can return as is
	return Kind(knd), instance
}

// Normalize normalizes this scope.
// It removes any uneeded parameters, and turns any invalid scopes into standard form.
func (scope Scope) Normalize() Scope {
	knd, instance := scope.Split()
	if instance == "" {
		return Scope(knd)
	}
	return Scope(string(knd) + ":" + instance)
}
