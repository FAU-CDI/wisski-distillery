package status

import (
	"encoding/json"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"golang.org/x/exp/slices"
)

// User represents a WissKI User
type User struct {
	UID     phpx.Integer   `json:"uid,omitempty"`
	Name    phpx.String    `json:"name,omitempty"`
	Mail    phpx.String    `json:"mail,omitempty"`
	Status  phpx.Boolean   `json:"status,omitempty"`
	Created phpx.Timestamp `json:"created,omitempty"`
	Changed phpx.Timestamp `json:"changed,omitempty"`
	Access  phpx.Timestamp `json:"access,omitempty"`
	Login   phpx.Timestamp `json:"login,omitempty"`
	Roles   UserRoles      `json:"roles,omitempty"`
}

// UserRole represents the role of a user
type UserRole string

const (
	Administrator UserRole = "administrator"
	ContentEditor UserRole = "content_editor"
)

// UserRoles represents a set of user roles for a given user
type UserRoles map[UserRole]struct{}

// Has checks if the UserRole has the given role
func (ur UserRoles) Has(role UserRole) (ok bool) {
	_, ok = ur[role]
	return
}

func (ur UserRoles) MarshalJSON() ([]byte, error) {
	roles := make([]string, len(ur))
	i := 0
	for r := range ur {
		roles[i] = string(r)
		i++
	}
	slices.Sort(roles) // for consistent marshaling

	return json.Marshal(strings.Join(roles, ", "))
}

func (u *UserRoles) UnmarshalJSON(data []byte) error {
	return phpx.UnmarshalIntermediate(u, func(s phpx.String) (UserRoles, error) {
		if len(s) == 0 {
			return nil, nil
		}
		roles := strings.Split(string(s), ", ")
		uroles := make(UserRoles, len(roles))
		for _, r := range roles {
			uroles[UserRole(r)] = struct{}{}
		}
		return uroles, nil
	}, data)
}
