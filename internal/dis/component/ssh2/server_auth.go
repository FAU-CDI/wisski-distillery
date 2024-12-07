package ssh2

//spellchecker:words github wisski distillery internal component sshkeys gliderlabs
import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2/sshkeys"
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

func (ssh2 *SSH2) handleAuth(ctx ssh.Context, key ssh.PublicKey) bool {
	return sshkeys.Slowdown(func() (ok bool) {
		permissions := make(map[string]bool)

		// grab the global permissions
		{
			globalKeys, err := ssh2.dependencies.Keys.Admin(ctx)
			if err != nil {
				return false
			}
			permissions[""] = sshkeys.KeyOneOf(globalKeys, key)
			ok = permissions[""]
		}

		// grab permissions for each instance
		{
			instances, err := ssh2.dependencies.Instances.All(ctx)
			if err != nil {
				return false
			}

			for _, instance := range instances {
				ikeys, err := instance.SSH().Keys(ctx)
				if err != nil {
					continue
				}
				access := sshkeys.KeyOneOf(ikeys, key)

				permissions[instance.Slug] = access || permissions[""]
				ok = ok || access
			}
		}

		setPermissions(ctx, permissions)
		return
	})
}
