//spellchecker:words triplestore
package triplestore

//spellchecker:words context errors http github wisski distillery internal component models pkglib errorsx
import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"go.tkw01536.de/pkglib/errorsx"
)

func (Triplestore) SnapshotNeedsRunning() bool { return false }

func (Triplestore) SnapshotName() string { return "triplestore" }

func (ts *Triplestore) Snapshot(wisski models.Instance, scontext *component.StagingContext) error {
	if err := scontext.AddDirectory(".", func(ctx context.Context) error {
		if err := scontext.AddFile(wisski.GraphDBRepository+".nq", func(ctx context.Context, file io.Writer) error {
			_, err := ts.SnapshotDB(ctx, file, wisski.GraphDBRepository)
			if err != nil {
				return fmt.Errorf("failed to snapshot database: %w", err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed to add nq file: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to add directory: %w", err)
	}
	return nil
}

var errTSBackupWrongStatusCode = errors.New("Triplestore.Backup: Wrong status code")

const nquadsContentType = "text/x-nquads"

// SnapshotDB snapshots the provided repository into dst.
func (ts Triplestore) SnapshotDB(ctx context.Context, dst io.Writer, repo string) (c int64, e error) {
	res, err := ts.DoRest(ctx, 0, http.MethodGet, "/repositories/"+url.PathEscape(repo)+"/statements?infer=false", &RequestHeaders{Accept: nquadsContentType})
	if err != nil {
		return 0, fmt.Errorf("failed to send rest request: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if res.StatusCode != http.StatusOK {
		return 0, errTSBackupWrongStatusCode
	}
	count, err := io.Copy(dst, res.Body)
	if err != nil {
		return count, fmt.Errorf("failed to copy result: %w", err)
	}
	return count, nil
}
