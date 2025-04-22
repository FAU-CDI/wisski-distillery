//spellchecker:words status
package status

//spellchecker:words time github wisski distillery internal config models
import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Distillery holds status and analytical data about a distillery.
type Distillery struct {
	Time time.Time // Time when this information was built

	// Configuration of the distillery
	Config *config.Config

	// number of instances
	TotalCount   int
	RunningCount int
	StoppedCount int

	Backups []models.Export // list of backups
}
