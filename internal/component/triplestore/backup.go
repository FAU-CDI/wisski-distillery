package triplestore

import (
	"encoding/json"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
)

func (ts *Triplestore) BackupName() string { return "triplestore" }

// Backup makes a backup of all Triplestore repositories databases into the path dest.
func (ts *Triplestore) Backup(context component.StagingContext) error {

	// list all the directories
	repos, err := ts.listRepositories()
	if err != nil {
		return err
	}

	// then backup each file separatly
	return context.AddDirectory("", func() error {
		for _, repo := range repos {
			if err := context.AddFile(repo.ID+".nq", func(file io.Writer) error {
				_, err := ts.SnapshotDB(file, repo.ID)
				return err
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (ts Triplestore) listRepositories() (repos []Repository, err error) {
	res, err := ts.OpenRaw("GET", "/rest/repositories", nil, "", "application/json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&repos)
	return
}
