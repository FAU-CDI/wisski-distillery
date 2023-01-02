package malt

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore"
)

// Malt is a component passed to every WissKI ingredient
type Malt struct {
	component.Base

	TS          *triplestore.Triplestore `auto:"true"`
	SQL         *sql.SQL                 `auto:"true"`
	Meta        *meta.Meta               `auto:"true"`
	ExporterLog *logger.Logger           `auto:"true"`
	Policy      *policy.Policy           `auto:"true"`
}
