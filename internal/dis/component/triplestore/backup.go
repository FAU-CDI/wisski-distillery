//spellchecker:words triplestore
package triplestore

//spellchecker:words context encoding json http github wisski distillery internal component pkglib errorsx
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

func (ts *Triplestore) BackupName() string { return "triplestore" }

// Backup makes a backup of all Triplestore repositories databases into the path dest.
func (ts *Triplestore) Backup(scontext *component.StagingContext) error {
	if err := scontext.AddDirectory("", func(ctx context.Context) error {
		// list all the directories
		repos, err := ts.client().ListRepositories(ctx)
		if err != nil {
			return fmt.Errorf("failed to list repositories: %w", err)
		}

		for _, repo := range repos {
			if err := scontext.AddFile(repo.ID+".nq", func(ctx context.Context, file io.Writer) error {
				_, err := ts.client().ExportContent(ctx, file, repo.ID)
				if err != nil {
					return fmt.Errorf("failed to snapshot database: %w", err)
				}
				return nil
			}); err != nil {
				return fmt.Errorf("failed to %s.nq: %w", repo.ID, err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to add directory: %w", err)
	}
	return nil
}
