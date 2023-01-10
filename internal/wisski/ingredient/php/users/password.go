package users

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/password"
)

var errGetValidator = errors.New("GetPasswordValidator: Unknown Error")

func (u *Users) GetPasswordValidator(ctx context.Context, username string) (pv PasswordValidator, err error) {
	server := u.Dependencies.PHP.NewServer()

	var hash string
	err = u.Dependencies.PHP.ExecScript(ctx, server, &hash, usersPHP, "get_password_hash", username)
	if err != nil {
		server.Close()
		return pv, err
	}
	if len(hash) == 0 {
		server.Close()
		return pv, errGetValidator
	}

	pv.server = server
	pv.username = username
	pv.hash = hash
	return pv, nil
}

type PasswordValidator struct {
	server *phpx.Server

	username string
	hash     string
}

func (pv PasswordValidator) Close() error {
	return pv.server.Close()
}

func (pv PasswordValidator) Check(ctx context.Context, password string) bool {
	var result phpx.Boolean
	err := pv.server.MarshalCall(ctx, &result, "check_password_hash", password, string(pv.hash))
	if err != nil {
		return false
	}
	return bool(result)
}

var errPasswordUsername = errors.New("username === password")

func (pv PasswordValidator) CheckDictionary(ctx context.Context, writer io.Writer) error {
	var counter int

	if pv.Check(ctx, pv.username) {
		if writer != nil {
			counter++
			fmt.Fprintln(writer, counter)
		}
		return errPasswordUsername
	}
	for candidate := range password.CommonPasswords() {
		if ctx.Err() != nil {
			continue
		}
		result := pv.Check(ctx, candidate.Password)
		if writer != nil {
			counter++
			fmt.Fprintln(writer, counter)
		}

		if result {
			return &password.CommonPasswordError{CommonPassword: candidate}
		}
	}

	return ctx.Err()
}
