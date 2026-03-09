//spellchecker:words triplestore
package triplestore

//spellchecker:words bytes context encoding json errors mime multipart http time github wisski distillery internal component wdlog pkglib errorsx timex
import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore/client"
)

// http.Client Timeout to be used for "trivial" triplestore operations.
// This includes e.g. CRUDing a specific repo.
const tsTrivialTimeout = time.Minute

// globalClient() returns a new globalClient for the triplestore API.
func (ts *Triplestore) globalClient() *client.Client {
	config := component.GetStill(ts).Config.TS
	client := client.NewClient(tsTrivialTimeout, ts.BaseURL, config.AdminUsername, config.AdminPassword)
	client.PollInterval = ts.PollInterval
	return client
}
