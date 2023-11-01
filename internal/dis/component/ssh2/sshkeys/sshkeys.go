package sshkeys

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/gliderlabs/ssh"
)

type SSHKeys struct {
	component.Base
	dependencies struct {
		SQL  *sql.SQL
		Auth *auth.Auth
	}
}

var (
	_ component.Table          = (*SSHKeys)(nil)
	_ component.UserDeleteHook = (*SSHKeys)(nil)
)

// Admin returns the set of administrative ssh keys.
// These are ssh keys associated to distillery admin users.
func (k *SSHKeys) Admin(ctx context.Context) (keys []ssh.PublicKey, err error) {
	users, err := k.dependencies.Auth.Users(ctx)
	if err != nil {
		return nil, err
	}

	// iterate over enabled distillery admin users
	for _, user := range users {
		if !user.IsEnabled() || !user.IsAdmin() {
			continue
		}
		ukeys, err := k.Keys(ctx, user.User.User)
		if err != nil {
			return nil, err
		}
		for _, ukey := range ukeys {
			if pk := ukey.PublicKey(); pk != nil {
				keys = append(keys, pk)
			}
		}
	}

	// and return the keys!
	return keys, nil
}
