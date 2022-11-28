package users

import (
	"bufio"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
)

var errGetValidator = errors.New("GetPasswordValidator: Unknown Error")

func (u *Users) GetPasswordValidator(ctx context.Context, username string) (pv PasswordValidator, err error) {
	server := u.PHP.NewServer()

	var hash string
	err = u.PHP.ExecScript(ctx, server, &hash, usersPHP, "get_password_hash", username)
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

type CommonPasswordError struct {
	Password CommonPassword
}

func (cpe CommonPasswordError) Error() string {
	return fmt.Sprintf("%q from %q", cpe.Password.Password, cpe.Password.Source)
}

func (pv PasswordValidator) CheckDictionary(ctx context.Context, writer io.Writer) error {
	var counter int

	if pv.Check(ctx, pv.username) {
		if writer != nil {
			counter++
			fmt.Fprintln(writer, counter)
		}
		return errPasswordUsername
	}
	for candidate := range CommonPasswords() {
		if ctx.Err() != nil {
			continue
		}
		result := pv.Check(ctx, candidate.Password)
		if writer != nil {
			counter++
			fmt.Fprintln(writer, counter)
		}

		if result {
			return &CommonPasswordError{Password: candidate}
		}
	}

	return ctx.Err()
}

//go:embed passwords
var passwordsEmbed embed.FS

type CommonPassword struct {
	Password string
	Source   string
}

// CommonPasswords returns a channel of most common passwords
func CommonPasswords() <-chan CommonPassword {
	pChan := make(chan CommonPassword, 10)
	go func() {
		defer close(pChan)

		fs.WalkDir(passwordsEmbed, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// get the full path
			if d.IsDir() || !strings.HasSuffix(path, ".txt") {
				return nil
			}

			// open it
			file, err := passwordsEmbed.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// scan it line by line
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "//") {
					continue
				}
				pChan <- CommonPassword{
					Password: line,
					Source:   path,
				}
			}

			return scanner.Err()
		})
	}()
	return pChan
}
