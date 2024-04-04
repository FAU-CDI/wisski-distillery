package triplestore

import (
	"context"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/pkg/errors"
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

// SnapshotDB snapshots the provided repository into dst
func (ts Triplestore) SnapshotDB(ctx context.Context, dst io.Writer, repo string) (int64, error) {
	res, err := ts.DoRest(ctx, 0, http.MethodGet, "/repositories/"+repo+"/statements?infer=false", &RequestHeaders{Accept: nquadsContentType})
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return 0, errTSBackupWrongStatusCode
	}
	return io.Copy(dst, res.Body)
}
