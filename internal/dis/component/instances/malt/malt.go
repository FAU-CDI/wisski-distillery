//spellchecker:words malt
package malt

//spellchecker:words github wisski distillery internal component auth policy docker exporter logger meta sshkeys triplestore
import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/docker"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2/sshkeys"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore"
)

// Malt is a component passed to every WissKI ingredient
type Malt struct {
	component.Base

	SQL           *sql.SQL           `inject:"true"`
	InstanceTable *sql.InstanceTable `inject:"true"`
	LockTable     *sql.LockTable     `inject:"true"`

	TS          *triplestore.Triplestore `inject:"true"`
	Meta        *meta.Meta               `inject:"true"`
	ExporterLog *logger.Logger           `inject:"true"`
	Policy      *policy.Policy           `inject:"true"`

	Docker *docker.Docker `inject:"true"`

	Keys *sshkeys.SSHKeys `inject:"true"`
}
