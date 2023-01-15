package ssh2

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2/sshkeys"
)

type SSH2 struct {
	component.Base
	Dependencies struct {
		SQL       *sql.SQL
		Instances *instances.Instances
		Auth      *auth.Auth
		Keys      *sshkeys.SSHKeys
	}
}

var (
	_ component.Installable = (*SSH2)(nil)
	_ component.Routeable   = (*SSH2)(nil)
)
