package triplestore

import (
	"context"
	"encoding/json"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

func (ts *Triplestore) BackupName() string { return "triplestore" }

// Backup makes a backup of all Triplestore repositories databases into the path dest.
func (ts *Triplestore) Backup(scontext *component.StagingContext) error {
	return scontext.AddDirectory("", func(ctx context.Context) error {
		// list all the directories
		repos, err := ts.listRepositories(ctx)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			if err := scontext.AddFile(repo.ID+".nq", func(ctx context.Context, file io.Writer) error {
				_, err := ts.SnapshotDB(ctx, file, repo.ID)
				return err
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (ts Triplestore) listRepositories(ctx context.Context) (repos []Repository, err error) {
	res, err := ts.OpenRaw(ctx, "GET", "/rest/repositories", nil, "", "application/json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&repos)
	return
}
