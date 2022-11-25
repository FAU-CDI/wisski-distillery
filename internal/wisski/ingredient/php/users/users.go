package users

import (
	_ "embed"
	"errors"
	"net/url"

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
func (u *Users) All(server *phpx.Server) (users []status.User, err error) {
	err = u.PHP.ExecScript(server, &users, usersPHP, "list_users")
	return
}

var errLoginUnknownError = errors.New("Login: Unknown Error")

// Login generates a login link for the user with the given username
func (u *Users) Login(server *phpx.Server, username string) (dest *url.URL, err error) {

	// generate a (relative) link
	var path string
	err = u.PHP.ExecScript(server, &path, usersPHP, "get_login_link", username)

	// if something went wrong, return
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, errLoginUnknownError
	}

	// parse it as a url
	dest, err = url.Parse(path)
	if err != nil {
		return nil, err
	}

	// and resolve the (possibly relative) reference
	dest = u.URL().ResolveReference(dest)
	return
}

var errSetPassword = errors.New("SetPassword: Unknown Error")

// SetPassword sets the password for a given user
func (u *Users) SetPassword(server *phpx.Server, username, password string) error {
	var ok bool
	err := u.PHP.ExecScript(server, &ok, usersPHP, "set_user_password", username, password)
	if err != nil {
		return err
	}
	if !ok {
		return errSetPassword
	}
	return nil
}

func (u *Users) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.Users, _ = u.All(flags.Server)
	return
}
