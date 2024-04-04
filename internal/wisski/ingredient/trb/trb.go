package trb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

type TRB struct {
	ingredient.Base

	dependencies struct {
		Barrel *barrel.Barrel
	}
}

func (trb *TRB) DoSomething(ctx context.Context, out io.Writer, allowEmptyRepository bool) (err error) {

	// stop instance, restart when done
	logging.LogMessage(out, "Shutting down instance")
	if err := trb.dependencies.Barrel.Stack().Down(ctx, out); err != nil {
		return err
	}

	defer func() {
		logging.LogMessage(out, "Restarting instance")
		e := trb.dependencies.Barrel.Stack().Up(ctx, out)
		if err == nil {
			err = e
		}
	}()

	// make the backup
	logging.LogMessage(out, "Dumping triplestore")
	path, err := trb.makeBackup(ctx, allowEmptyRepository)
	if err != nil {
		return err
	}
	fmt.Printf("Wrote %q\n", path)

	logging.LogMessage(out, "Purging triplestore")
	if err := trb.Malt.TS.Purge(ctx, trb.Instance, trb.Domain()); err != nil {
		return err
	}

	logging.LogMessage(out, "Provising triplestore")
	if err := trb.Malt.TS.Provision(ctx, trb.Instance, trb.Domain()); err != nil {
		return err
	}

	logging.LogMessage(out, "Loading dump file")
	content, err := os.Open(path)
	if err != nil {
		return err
	}
	defer content.Close()

	logging.LogMessage(out, "Restoring triplestore")
	if err := trb.Malt.TS.RestoreDB(ctx, trb.GraphDBRepository, content); err != nil {
		return err
	}

	return
}

var errBackupEmpty = errors.New("no data contained in backup file (is the repository empty?)")

func (trb *TRB) makeBackup(ctx context.Context, allowEmptyRepository bool) (path string, err error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	defer f.Close()

	count, err := trb.Malt.TS.SnapshotDB(ctx, f, trb.GraphDBRepository)
	if err != nil {
		return "", err
	}

	if count == 0 && !allowEmptyRepository {
		return "", errBackupEmpty
	}

	return f.Name(), nil
}
