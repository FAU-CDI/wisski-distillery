package password

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

// CommonPasswordError
type CommonPasswordError struct {
	CommonPassword
}

func (cpe CommonPasswordError) Error() string {
	return fmt.Sprintf("%q from %q", cpe.Password, cpe.Source)
}

type CommonPassword struct {
	Password string
	Source   string
}

//go:embed common
var commonEmbed embed.FS

// CommonPasswords returns a channel that contains all passwords.
// The caller must drain the channel.
func CommonPasswords() <-chan CommonPassword {
	pChan := make(chan CommonPassword, 10)
	go func() {
		defer close(pChan)

		fs.WalkDir(commonEmbed, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// get the full path
			if d.IsDir() || !strings.HasSuffix(path, ".txt") {
				return nil
			}

			// open it
			file, err := commonEmbed.Open(path)
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

// CheckCommonPassword checks if a password is a common password.
//
// check is called with each candidate password to perform the check.
// check should return a boolean indicating if the password in question corresponds to the candidate.
//
// CheckCommonPassword returns one of three error values.
//
// - a CommonPasswordError (when a password matches a common password)
// - an error returned by check (assuming some check went wrong)
// - or nil (when a password is not a common password
func CheckCommonPassword(check func(candidate string) (bool, error)) error {
	for commmon := range CommonPasswords() {
		ok, err := check(commmon.Password)
		if err != nil {
			return err
		}

		// password validation passed
		if ok {
			return CommonPasswordError{
				CommonPassword: commmon,
			}
		}
	}
	return nil
}
