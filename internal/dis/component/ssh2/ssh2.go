package ssh2

//spellchecker:words github wisski distillery internal component auth instances sshkeys pkglib lazy
import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2/sshkeys"
	"github.com/tkw1536/pkglib/lazy"
)

type SSH2 struct {
	component.Base
	dependencies struct {
		SQL       *sql.SQL
		Instances *instances.Instances
		Auth      *auth.Auth
		Keys      *sshkeys.SSHKeys
	}

	interceptsC lazy.Lazy[[]Intercept]
}

var (
	_ component.Installable = (*SSH2)(nil)
	_ component.Routeable   = (*SSH2)(nil)
)
