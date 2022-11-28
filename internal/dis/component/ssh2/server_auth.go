package ssh2

import (
	"time"

	"github.com/gliderlabs/ssh"
)

func (ssh2 *SSH2) setupAuth(server *ssh.Server) {
	server.PublicKeyHandler = ssh2.handleAuth
}

// ssh2Key is a type of context keys for this package
type ssh2Key int

const (
	// permissions represents the permissions for the given session
	permission ssh2Key = iota
)

func setPermissions(context ssh.Context, permissions map[string]bool) {
	context.SetValue(permission, permissions)
}

// hasPermission checks if the given context permits access to the given slug.
// The empty slug checks for global access.
func hasPermission(context ssh.Context, slug string) bool {
	value, ok := context.Value(permission).(map[string]bool)
	return ok && value[slug]
}

// getAnyPermission gets some instance the user has access to.
// If the user does not have access to anything, returns "", false.
// If the user has superuser access, but there are no instances, returns "", true.
func getAnyPermission(context ssh.Context) (string, bool) {
	value, ok := context.Value(permission).(map[string]bool)
	if !ok {
		return "", false
	}

	for slug, ok := range value {
		if ok && slug != "" {
			return slug, true
		}
	}

	return "", (false || value[""])
}

const authDelay = time.Second / 10

func (ssh2 *SSH2) handleAuth(ctx ssh.Context, key ssh.PublicKey) bool {
	return slowdown(func() (ok bool) {
		permissions := make(map[string]bool)

		// grab the global permissions
		{
			globalKeys, err := ssh2.GlobalKeys()
			if err != nil {
				return false
			}
			permissions[""] = isKey(globalKeys, key)
			ok = permissions[""]
		}

		// grab permissions for each instance
		{
			instances, err := ssh2.Instances.All(ctx)
			if err != nil {
				return false
			}

			for _, instance := range instances {
				ikeys, err := instance.SSH().Keys()
				if err != nil {
					continue
				}
				access := isKey(ikeys, key)

				permissions[instance.Slug] = access || permissions[""]
				ok = ok || access
			}
		}

		setPermissions(ctx, permissions)
		return
	}, authDelay)
}

// slowdown invokes f immediatly, but only returns the result to the caller after at least duration.
// It can be used to prevent timing attacks
func slowdown[T any](f func() T, duration time.Duration) T {
	result := make(chan T, 1)
	go func() {
		result <- f()
	}()
	time.Sleep(duration)
	return <-result
}

// isKey checks if keys contains key in O(len(keys))
func isKey(keys []ssh.PublicKey, key ssh.PublicKey) bool {
	var res bool
	for _, ak := range keys {
		if ssh.KeysEqual(ak, key) {
			res = true
		}
	}
	return res
}
