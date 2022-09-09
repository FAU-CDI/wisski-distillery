package cmd

import (
	"io/fs"
	"os"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/FAU-CDI/wisski-distillery/internal/targz"
	"github.com/tkw1536/goprogram/exit"
)

// Snapshot creates a snapshot of an instance
var Snapshot wisski_distillery.Command = snapshot{}

type snapshot struct {
	Keepalive   bool `short:"k" long:"keepalive" description:"Keep instance running while taking a backup. Might lead to inconsistent state"`
	StagingOnly bool `short:"s" long:"staging-only" description:"Do not package into a snapshot archive, but only create a staging directory"`

	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to take a snapshot of"`
		Dest string `positional-arg-name:"DEST" description:"Destination path to write snapshot archive to. Defaults to the snapshots/archives/ directory"`
	} `positional-args:"true"`
}

func (snapshot) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: env.Requirements{
			NeedsDistillery: true,
		},
		Command:     "snapshot",
		Description: "Generates a snapshot archive for the provided archive",
	}
}

var errSnapshotFailed = exit.Error{
	Message:  "Failed to make a snapshot",
	ExitCode: exit.ExitGeneric,
}

func (bi snapshot) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	instance, err := dis.Instance(bi.Positionals.Slug)
	if err != nil {
		return err
	}

	logging.LogMessage(context.IOStream, "Creating snapshot of instance %s", bi.Positionals.Slug)

	// determine the target path for the archive
	var sPath string
	if !bi.StagingOnly {
		// regular mode: create a temporary staging directory
		logging.LogMessage(context.IOStream, "Creating new snapshot staging directory")
		sPath, err = dis.NewSnapshotStagingDir(instance.Slug)
		if err != nil {
			return errSnapshotFailed.Wrap(err)
		}
		defer func() {
			logging.LogMessage(context.IOStream, "Removing snapshot staging directory")
			os.RemoveAll(sPath)
		}()
	} else {
		// staging mode: use dest as a destination
		sPath = bi.Positionals.Dest
		if sPath == "" {
			sPath, err = dis.NewSnapshotStagingDir(instance.Slug)
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

	// TODO: Allow skipping backups of individual parts and make them concurrent!

	// take a snapshot into the staging area!
	logging.LogOperation(func() error {
		sreport := instance.Snapshot(context.IOStream, env.SnapshotDescription{
			Dest:      sPath,
			Keepalive: bi.Keepalive,
		})

		// write out the report, ignoring any errors!
		sreport.WriteReport(context.IOStream)

		return nil
	}, context.IOStream, "Generating Snapshot")

	// if we requested to only have a staging area, then we are done
	if bi.StagingOnly {
		context.Printf("Wrote %s\n", sPath)
		return nil
	}

	// create the archive path
	archivePath := bi.Positionals.Dest
	if archivePath == "" {
		archivePath = dis.NewSnapshotArchivePath(instance.Slug)
	}

	// and write everything into it!
	// TODO: Should we move the open call to here?
	var count int64
	if err := logging.LogOperation(func() error {
		context.IOStream.Println(archivePath)

		count, err = targz.Package(archivePath, sPath, func(dst, src string) {
			context.Println(dst)
		})
		return err
	}, context.IOStream, "Writing snapshot archive"); err != nil {
		return errSnapshotFailed.Wrap(err)
	}
	context.Printf("Wrote %d byte(s) to %s\n", count, archivePath)
	return nil
}
