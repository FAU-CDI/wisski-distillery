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
	logging.LogMessage(out, "Re-creating adapter")
	if _, err := trb.dependencies.Adapters.SetAdapter(ctx, nil, trb.dependencies.Adapters.DefaultAdapter()); err != nil {
		return 0, err
	}

	// stop instance, restart when done
	logging.LogMessage(out, "Shutting down instance")
	if err := trb.dependencies.Barrel.Stack().Down(ctx, out); err != nil {
		return 0, err
	}

	defer func() {
		logging.LogMessage(out, "Restarting instance")
		e := trb.dependencies.Barrel.Stack().Up(ctx, out)
		if err == nil {
			err = e
		}
	}()

	// make the backup
	logging.LogMessage(out, "Storing triplestore content")
	dumpPath, _, err := trb.makeBackup(ctx, allowEmptyRepository)
	if err != nil {
		return 0, err
	}
	fmt.Printf("Wrote %q\n", dumpPath)

	liquid := ingredient.GetLiquid(trb)

	logging.LogMessage(out, "Purging triplestore")
	if err := liquid.TS.Purge(ctx, liquid.Instance, liquid.Domain()); err != nil {
		return 0, err
	}

	logging.LogMessage(out, "Provising triplestore")
	if err := liquid.TS.Provision(ctx, liquid.Instance, liquid.Domain()); err != nil {
		return 0, err
	}

	logging.LogMessage(out, "Restoring triplestore")
	if err := trb.restoreBackup(ctx, dumpPath); err != nil {
		return 0, err
	}

	logging.LogMessage(out, "Deleting dump file")
	if err := os.Remove(dumpPath); err != nil {
		return 0, err
	}

	return
}

var errBackupEmpty = errors.New("no data contained in backup file (is the repository empty?)")

func (trb *TRB) makeBackup(ctx context.Context, allowEmptyRepository bool) (path string, size int64, err error) {
	file, err := os.CreateTemp("", "*.nq.gz")
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	// create a new writer
	zippedFile := gzip.NewWriter(file)
	defer zippedFile.Close()

	{
		liquid := ingredient.GetLiquid(trb)
		size, err := liquid.TS.SnapshotDB(ctx, zippedFile, liquid.GraphDBRepository)
		if err != nil {
			return "", 0, err
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
		return err
	}
	defer reader.Close()

	decompressedReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer decompressedReader.Close()

	liquid := ingredient.GetLiquid(trb)
	if err := liquid.TS.RestoreDB(ctx, liquid.GraphDBRepository, decompressedReader); err != nil {
		return err
	}
	return nil
}
