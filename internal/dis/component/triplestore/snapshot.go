//spellchecker:words triplestore
package triplestore

//spellchecker:words context http github wisski distillery internal component models errors
import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

func (Triplestore) SnapshotNeedsRunning() bool { return false }

func (Triplestore) SnapshotName() string { return "triplestore" }

func (ts *Triplestore) Snapshot(wisski models.Instance, scontext *component.StagingContext) error {
	return scontext.AddDirectory(".", func(ctx context.Context) error {
		return scontext.AddFile(wisski.GraphDBRepository+".nq", func(ctx context.Context, file io.Writer) error {
			_, err := ts.SnapshotDB(ctx, file, wisski.GraphDBRepository)
			return err
		})
	})
}

var errTSBackupWrongStatusCode = errors.New("Triplestore.Backup: Wrong status code")

const nquadsContentType = "text/x-nquads"

// SnapshotDB snapshots the provided repository into dst.
func (ts Triplestore) SnapshotDB(ctx context.Context, dst io.Writer, repo string) (c int64, e error) {
	res, err := ts.DoRest(ctx, 0, http.MethodGet, "/repositories/"+repo+"/statements?infer=false", &RequestHeaders{Accept: nquadsContentType})
	if err != nil {
		return 0, fmt.Errorf("failed to send rest request: %w", err)
	}
	defer func() {
		e2 := res.Body.Close()
		if e2 == nil {
			return
		}
		e2 = fmt.Errorf("failed to close body: %w", e2)
		if e == nil {
			e = e2
		} else {
			e = errors.Join(e, e2)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return 0, errTSBackupWrongStatusCode
	}
	return io.Copy(dst, res.Body)
}
