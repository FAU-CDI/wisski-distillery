package triplestore

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tkw1536/goprogram/stream"
)

func (ts *Triplestore) BackupName() string { return "triplestore" }

// Backup makes a backup of all Triplestore repositories databases into the path dest.
func (ts *Triplestore) Backup(io stream.IOStream, dest string) error {

	// list all the repositories
	repos, err := ts.listRepositories()
	if err != nil {
		return err
	}

	// create the base directory, todo: outsource this
	if err := os.Mkdir(dest, fs.ModeDir); err != nil {
		return err
	}

	// iterate over all the repositories
	for _, repo := range repos {
		if rErr := (func(repo Repository) error {
			name := filepath.Join(dest, repo.ID+".nq")

			// todo: outsource this
			dest, err := os.Create(name)
			if err != nil {
				return err
			}
			defer dest.Close()

			_, err = ts.Snapshot(dest, repo.ID)
			return err
		}(repo)); err == nil && rErr != nil {
			err = rErr
		}
	}
	return err
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
