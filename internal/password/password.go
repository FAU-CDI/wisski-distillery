// Package password allows generating random passwords
package password

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// NOTE(twiesing): A bunch of scripts cannot properly handle the extra characters in the password.
// For now it is disabled, but it should be re-enabled later.
const PasswordCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // + "!@#$%&*"
const PasswordCharCount = len(PasswordCharSet)

// Password returns a randomly generated password with the provided length.
// [rand.Reader] is used as the source of randomness.
func Password(length int) (string, error) {
	if length < 0 {
		panic("length < 0")
	}

	var password strings.Builder
	password.Grow(length)

	for i := 0; i < length; i++ {

		// grab a random index!
		index, err := rand.Int(rand.Reader, big.NewInt(int64(PasswordCharCount)))
		if err != nil {
			return "", err
		}

		// and use that index!
		if err := password.WriteByte(PasswordCharSet[int(index.Int64())]); err != nil {
			return "", err
		}
	}

	// return the password!
	return password.String(), nil
}
