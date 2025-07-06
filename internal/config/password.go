//spellchecker:words config
package config

//spellchecker:words crypto rand github wisski distillery internal passwordx pkglib password
import (
	"crypto/rand"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/passwordx"
	"go.tkw01536.de/pkglib/password"
)

// NewPassword returns a new password using the password settings from this configuration.
func (cfg Config) NewPassword() (string, error) {
	pass, err := password.Generate(rand.Reader, cfg.PasswordLength, passwordx.Safe)
	if err != nil {
		return pass, fmt.Errorf("failed to generate password: %w", err)
	}
	return pass, nil
}
