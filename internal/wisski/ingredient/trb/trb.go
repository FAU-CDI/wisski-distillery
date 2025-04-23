package trb

//spellchecker:words compress gzip context errors github wisski distillery internal ingredient barrel logging
import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

type TRB struct {
	ingredient.Base

	dependencies struct {
		Barrel   *barrel.Barrel
		Adapters *extras.Adapters
	}
}

// RebuildTriplestore rebuilds the triplestore by making a backup, storing it on disk, purging the triplestore, and restoring the backup.
// Returns the size of the backup dump in bytes.
func (trb *TRB) RebuildTriplestore(ctx context.Context, out io.Writer, allowEmptyRepository bool) (size int, err error) {
	// re-create the default adapter
	if _, err := logging.LogMessage(out, "Re-creating adapter"); err != nil {
		return 0, fmt.Errorf("failed to log message: %w", err)
	}
	if _, err := trb.dependencies.Adapters.SetAdapter(ctx, nil, trb.dependencies.Adapters.DefaultAdapter()); err != nil {
		return 0, err
	}

	// stop instance, restart when done
	if _, err := logging.LogMessage(out, "Shutting down instance"); err != nil {
		return 0, fmt.Errorf("failed to log message: %w", err)
	}
	if err := trb.dependencies.Barrel.Stack().Down(ctx, out); err != nil {
		return 0, err
	}

	defer func() {
		if _, e := logging.LogMessage(out, "Restarting instance"); e != nil {
			e = fmt.Errorf("failed to log message: %w", err)
			if err == nil {
				err = e
			} else {
				err = errors.Join(err, e)
			}
			return
		}

		e2 := trb.dependencies.Barrel.Stack().Up(ctx, out)
		if e2 == nil {
			return
		}
		if err == nil {
			err = e2
		} else {
			err = errors.Join(err, e2)
		}
	}()

	// make the backup
	if _, err := logging.LogMessage(out, "Storing triplestore content"); err != nil {
		return 0, fmt.Errorf("failed to log message: %w", err)
	}
	dumpPath, _, err := trb.makeBackup(ctx, allowEmptyRepository)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Wrote %q\n", dumpPath)

	liquid := ingredient.GetLiquid(trb)

	if _, err := logging.LogMessage(out, "Purging triplestore"); err != nil {
		return 0, fmt.Errorf("failed to log message: %w", err)
	}
	if err := liquid.TS.Purge(ctx, liquid.Instance, liquid.Domain()); err != nil {
		return 0, err
	}

	if _, err := logging.LogMessage(out, "Provising triplestore"); err != nil {
		return 0, fmt.Errorf("failed to log message: %w", err)
	}
	if err := liquid.TS.Provision(ctx, liquid.Instance, liquid.Domain()); err != nil {
		return 0, err
	}

	if _, err := logging.LogMessage(out, "Restoring triplestore"); err != nil {
		return 0, fmt.Errorf("failed to log message: %w", err)
	}
	if err := trb.restoreBackup(ctx, dumpPath); err != nil {
		return 0, err
	}

	if _, err := logging.LogMessage(out, "Deleting dump file"); err != nil {
		return 0, fmt.Errorf("failed to log message: %w", err)
	}
	if err := os.Remove(dumpPath); err != nil {
		return 0, fmt.Errorf("failed to delete dump file: %w", err)
	}

	return
}

var errBackupEmpty = errors.New("no data contained in backup file (is the repository empty?)")

func (trb *TRB) makeBackup(ctx context.Context, allowEmptyRepository bool) (path string, size int64, e error) {
	file, err := os.CreateTemp("", "*.nq.gz")
	if err != nil {
		return "", 0, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer func() {
		e2 := file.Close()
		if e2 == nil {
			return
		}
		e2 = fmt.Errorf("failed to close file: %w", e2)
		if e == nil {
			e = e2
		} else {
			e = errors.Join(e, e2)
		}
	}()

	// create a new writer
	zippedFile := gzip.NewWriter(file)
	defer zippedFile.Close()

	{
		liquid := ingredient.GetLiquid(trb)
		size, err := liquid.TS.SnapshotDB(ctx, zippedFile, liquid.GraphDBRepository)
		if err != nil {
			return "", 0, fmt.Errorf("failed to snapshot db: %w", err)
		}

		if size == 0 && !allowEmptyRepository {
			return "", 0, errBackupEmpty
		}

		return file.Name(), size, nil
	}
}

func (trb *TRB) restoreBackup(ctx context.Context, path string) (err error) {
	reader, err := os.Open(path) // #nosec G304 -- intended
	if err != nil {
		return fmt.Errorf("failed to restore database: %w", err)
	}
	defer reader.Close()

	decompressedReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer decompressedReader.Close()

	liquid := ingredient.GetLiquid(trb)
	if err := liquid.TS.RestoreDB(ctx, liquid.GraphDBRepository, decompressedReader); err != nil {
		return fmt.Errorf("failed to restore database: %w", err)
	}
	return nil
}
