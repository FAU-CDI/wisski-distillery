package triplestore

import (
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/pkg/errors"
)

func (Triplestore) SnapshotNeedsRunning() bool { return false }

func (Triplestore) SnapshotName() string { return "triplestore" }

func (ts *Triplestore) Snapshot(wisski models.Instance, context component.StagingContext) error {
	return context.AddDirectory(".", func() error {
		return context.AddFile(wisski.GraphDBRepository+".nq", func(file io.Writer) error {
			_, err := ts.SnapshotDB(file, wisski.GraphDBRepository)
			return err
		})
	})
}

var errTSBackupWrongStatusCode = errors.New("Triplestore.Backup: Wrong status code")

// SnapshotDB snapshots the provided repository into dst
func (ts Triplestore) SnapshotDB(dst io.Writer, repo string) (int64, error) {
	res, err := ts.OpenRaw("GET", "/repositories/"+repo+"/statements?infer=false", nil, "", "application/n-quads")
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, errTSBackupWrongStatusCode
	}
	defer res.Body.Close()
	return io.Copy(dst, res.Body)
}
