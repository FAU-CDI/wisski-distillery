//spellchecker:words status
package status

//spellchecker:words encoding json strings time slices github wisski distillery internal phpx
import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"slices"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
)

// DrupalUser represents a WissKI DrupalUser.
type DrupalUser struct {
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

func (du DrupalUser) String() string {
	var builder strings.Builder

	builder.WriteString("DrupalUser{")
	defer builder.WriteString("}")

	fmt.Fprintf(&builder, "UID: %d, ", du.UID)
	fmt.Fprintf(&builder, "Name: %q, ", du.Name)

	if du.Mail != "" {
		fmt.Fprintf(&builder, "Mail: %q, ", du.Mail)
	}

	fmt.Fprintf(&builder, "Status: %t, ", du.Status)

	for _, tn := range []struct {
		Name string
		Time time.Time
	}{
		{"Created", du.Created.Time()},
		{"Changed", du.Changed.Time()},
		{"Access", du.Access.Time()},
		{"Login", du.Login.Time()},
	} {
		if tn.Time.IsZero() {
			continue
		}
		fmt.Fprintf(&builder, "%s: %q, ", tn.Name, tn.Time.Format(time.Stamp))
	}

	fmt.Fprintf(&builder, "Roles: %s", du.Roles)

	return builder.String()
}

// UserRole represents the role of a user.
type UserRole string

const (
	Administrator UserRole = "administrator"
	ContentEditor UserRole = "content_editor"
)

// UserRoles represents a set of user roles for a given user.
//
//nolint:recvcheck
type UserRoles map[UserRole]struct{}

func (ur UserRoles) String() string {
	return "[" + ur.string() + "]"
}

// String turns this UserRoles into a string.
func (ur UserRoles) string() string {
	roles := make([]string, len(ur))
	i := 0
	for r := range ur {
		roles[i] = string(r)
		i++
	}
	slices.Sort(roles) // for consistent marshaling
	return strings.Join(roles, ", ")
}

func (ur UserRoles) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(ur.string())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal string: %w", err)
	}
	return bytes, nil
}

// Has checks if the UserRole has the given role.
func (ur UserRoles) Has(role UserRole) (ok bool) {
	_, ok = ur[role]
	return
}

func (u *UserRoles) UnmarshalJSON(data []byte) error {
	if err := phpx.UnmarshalIntermediate(u, func(s phpx.String) (UserRoles, error) {
		if len(s) == 0 {
			return nil, nil
		}
		roles := strings.Split(string(s), ", ")
		uroles := make(UserRoles, len(roles))
		for _, r := range roles {
			uroles[UserRole(r)] = struct{}{}
		}
		return uroles, nil
	}, data); err != nil {
		return fmt.Errorf("failed to unmarshal user roles: %w", err)
	}
	return nil
}
