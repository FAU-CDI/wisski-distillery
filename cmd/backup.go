package cmd

import (
	"io/fs"
	"os"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/backup"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/targz"
	"github.com/tkw1536/goprogram/exit"
)

// Backup is the 'backup' command
var Backup wisski_distillery.Command = backupC{}

type backupC struct {
	NoPrune     bool `short:"n" long:"no-prune" description:"Do not prune older backup archives"`
	StagingOnly bool `short:"s" long:"staging-only" description:"Do not package into a backup archive, but only create a staging directory"`
	Positionals struct {
		Dest string `positional-arg-name:"DEST" description:"Destination path to write backup archive to. Defaults to the snapshots/archives/ directory"`
	} `positional-args:"true"`
}

func (backupC) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "backup",
		Description: "Makes a backup of the entire distillery",
	}
}

var errBackupFailed = exit.Error{
	Message:  "Failed to make a backup",
	ExitCode: exit.ExitGeneric,
}

func (bk backupC) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	var err error

	if !bk.NoPrune {
		defer logging.LogOperation(func() error {
			return dis.PruneBackups(context.IOStream)
		}, context.IOStream, "Pruning old backups")
	}

	// determine the target path for the archive
	var sPath string
	if !bk.StagingOnly {
		// regular mode: create a temporary staging directory
		logging.LogMessage(context.IOStream, "Creating new snapshot staging directory")
		sPath, err = dis.NewSnapshotStagingDir("")
		if err != nil {
			return errSnapshotFailed.Wrap(err)
		}
		defer func() {
			logging.LogMessage(context.IOStream, "Removing snapshot staging directory")
			os.RemoveAll(sPath)
		}()
	} else {
		// staging mode: use dest as a destination
		sPath = bk.Positionals.Dest
		if sPath == "" {
			sPath, err = dis.NewSnapshotStagingDir("")
			if err != nil {
				return errSnapshotFailed.Wrap(err)
			}
		}

		// create the directory (if it doesn't already exist)
		logging.LogMessage(context.IOStream, "Creating staging directory")
		err = os.Mkdir(sPath, fs.ModePerm)
		if !os.IsExist(err) && err != nil {
			return errSnapshotFailed.WithMessageF(err)
		}
		err = nil
	}
	context.Println(sPath)

	logging.LogOperation(func() error {
		backup := backup.New(context.IOStream, dis, backup.Description{
			Dest: sPath,
			Auto: bk.Positionals.Dest == "",
		})
		backup.WriteReport(context.IOStream)
		return nil
	}, context.IOStream, "Generating Backup")

	// if we requested to only have a staging area, then we are done
	if bk.StagingOnly {
		context.Printf("Wrote %s\n", sPath)
		return nil
	}

	// create the archive path
	archivePath := bk.Positionals.Dest
	if archivePath == "" {
		archivePath = dis.NewSnapshotArchivePath("")
	}

	// and write everything into it!
	// TODO: Should we move the open call to here?
	var count int64
	if err := logging.LogOperation(func() error {
		context.IOStream.Println(archivePath)

		count, err = targz.Package(archivePath, sPath, func(dst, src string) {
			context.Printf("\033[2K\r%s", dst)
		})
		context.Println("")
		return err
	}, context.IOStream, "Writing backup archive"); err != nil {
		return errSnapshotFailed.Wrap(err)
	}
	context.Printf("Wrote %d byte(s) to %s\n", count, archivePath)

	return nil
}
