package malt

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2/sshkeys"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore"
)

// Malt is a component passed to every WissKI ingredient
type Malt struct {
	component.Base

	SQL           *sql.SQL           `auto:"true"`
	InstanceTable *sql.InstanceTable `auto:"true"`
	LockTable     *sql.LockTable     `auto:"true"`

	TS          *triplestore.Triplestore `auto:"true"`
	Meta        *meta.Meta               `auto:"true"`
	ExporterLog *logger.Logger           `auto:"true"`
	Policy      *policy.Policy           `auto:"true"`

	Keys *sshkeys.SSHKeys `auto:"true"`
}
