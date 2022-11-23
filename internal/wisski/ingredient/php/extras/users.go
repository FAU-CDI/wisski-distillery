package extras

import (
	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
)

type Users struct {
	ingredient.Base

	PHP *php.PHP
}

//go:embed users.php
var usersPHP string

// All returns all known usernames
func (u *Users) All(server *phpx.Server) (users []string, err error) {
	err = u.PHP.ExecScript(server, &users, usersPHP, "list_users")
	return
}

func (u *Users) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.Users, _ = u.All(flags.Server)
	return
}
