//spellchecker:words sshkeys
package sshkeys

//spellchecker:words crypto rand math time github gliderlabs
import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/gliderlabs/ssh"
)

// KeyOneOf checks if keys is one of the given set of keys.
func KeyOneOf(keys []ssh.PublicKey, key ssh.PublicKey) bool {
	return len(KeyIndexes(keys, key)) > 0
}

// KeyIndexes returns a slice of ints that contain the indexes of the given key.
func KeyIndexes(keys []ssh.PublicKey, key ssh.PublicKey) []int {
	indexes := make([]int, 0, len(keys))
	for i, cey := range keys {
		if ssh.KeysEqual(key, cey) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

const (
	slowdownMinDelay = time.Second / 10
	slowdownJitter   = time.Second / 10
)

// slowdown invokes f immediatly, but introduces a random delay to prevent timing attacks.
// the delay is also introduced if f() panics.
func Slowdown[T any](f func() T) T {
	start := time.Now()
	defer func() {
		// sleep the minimum remaining time
		remain := time.Since(start) - slowdownMinDelay
		if remain > 0 {
			time.Sleep(remain)
		}

		// find a second random delay
		delay, err := rand.Int(rand.Reader, big.NewInt(int64(slowdownJitter)))
		if err != nil {
			return
		}

		// and wait that long
		time.Sleep(time.Duration(delay.Int64()))
	}()

	return f()
}
