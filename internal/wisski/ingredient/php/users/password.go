//spellchecker:words users
package users

//spellchecker:words context errors github wisski distillery internal passwordx phpx pkglib errorsx password
import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/passwordx"
	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/password"
)

var errGetValidator = errors.New("`GetPasswordValidator': unknown error")

func (u *Users) GetPasswordValidator(ctx context.Context, username string) (pv PasswordValidator, err error) {
	server := u.dependencies.PHP.NewServer()

	var hash string
	err = u.dependencies.PHP.ExecScript(ctx, server, &hash, usersPHP, "get_password_hash", username)
	if err != nil {
		if e2 := server.Close(); e2 != nil {
			err = errors.Join(
				fmt.Errorf("failed to get password hash: %w", err),
				fmt.Errorf("failed to close server: %w", err),
			)
		}
		return pv, err
	}
	if len(hash) == 0 {
		return pv, errorsx.Combine(
			errGetValidator,
			server.Close(),
		)
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
	if err := pv.server.Close(); err != nil {
		return fmt.Errorf("failed to close php server: %w", err)
	}
	return nil
}

func (pv PasswordValidator) Check(ctx context.Context, password string) bool {
	var result phpx.Boolean
	err := pv.server.MarshalCall(ctx, &result, "check_password_hash", password, pv.hash)
	if err != nil {
		// TODO: Log?
		return false
	}
	return bool(result)
}

var errPasswordUsername = errors.New("username equals password")

func (pv PasswordValidator) CheckDictionary(ctx context.Context, writer io.Writer) error {
	var counter int

	if pv.Check(ctx, pv.username) {
		if writer != nil {
			counter++
			if _, err := fmt.Fprintln(writer, counter); err != nil {
				return fmt.Errorf("unable to report progress: %w", err)
			}
		}
		return errPasswordUsername
	}
	for candidate := range password.Passwords(passwordx.Sources...) {
		if ctx.Err() != nil {
			continue
		}
		result := pv.Check(ctx, candidate.Password)
		if writer != nil {
			counter++
			if _, err := fmt.Fprintln(writer, counter); err != nil {
				return fmt.Errorf("unable to report progress: %w", err)
			}
		}

		if result {
			return &password.CommonPasswordError{CommonPassword: candidate}
		}
	}

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context closed before returning: %w", err)
	}
	return nil
}
