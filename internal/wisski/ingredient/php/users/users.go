package users

import (
	"context"
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
	dependencies struct {
		PHP *php.PHP
	}
}

var (
	_ ingredient.WissKIFetcher = (*Users)(nil)
)

//go:embed users.php
var usersPHP string

// All returns all known usernames
func (u *Users) All(ctx context.Context, server *phpx.Server) (users []status.DrupalUser, err error) {
	err = u.dependencies.PHP.ExecScript(ctx, server, &users, usersPHP, "list_users")
	return
}

var errLoginUnknownError = errors.New("`Login': unknown error")

// Login generates a login link for the user with the given username
func (u *Users) Login(ctx context.Context, server *phpx.Server, username string) (dest *url.URL, err error) {
	return u.LoginWithOpt(ctx, server, username, LoginOptions{
		Destination:     "/",
		CreateIfMissing: false,
		GrantAdminRole:  false,
	})
}

type LoginOptions struct {
	Destination     string
	CreateIfMissing bool
	GrantAdminRole  bool
}

// LoginOrCreate generates a login link for the user with the given username and options
func (u *Users) LoginWithOpt(ctx context.Context, server *phpx.Server, username string, opts LoginOptions) (dest *url.URL, err error) {

	// generate a (relative) link
	var path string
	err = u.dependencies.PHP.ExecScript(ctx, server, &path, usersPHP, "get_login_link", username, opts.Destination, opts.CreateIfMissing, opts.GrantAdminRole)

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

var errSetPassword = errors.New("`SetPassword': unknown error")

// SetPassword sets the password for a given user
func (u *Users) SetPassword(ctx context.Context, server *phpx.Server, username, password string) error {
	var ok bool
	err := u.dependencies.PHP.ExecScript(ctx, server, &ok, usersPHP, "set_user_password", username, password)
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

	info.Users, _ = u.All(flags.Context, flags.Server)
	return
}
