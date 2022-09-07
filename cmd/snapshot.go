package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/FAU-CDI/wisski-distillery/internal/targz"
	"github.com/tkw1536/goprogram/exit"
)

// Snapshot creates a snapshot of an instance
var Snapshot wisski_distillery.Command = snapshot{}

type snapshot struct {
	Keepalive bool `short:"k" long:"keepalive" description:"Keep instance running while taking a backup. Might lead to inconsistent state"`

	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to take a snapshot of"`
		Dest string `positional-arg-name:"DEST" description:"Destination path to write snapshot archive to. Defaults to the snapshots/archives/ directory"`
	} `positional-args:"true"`
}

func (snapshot) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: env.Requirements{
			NeedsConfig: true,
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

	// TODO: Allow skipping backups of individual parts and make them concurrent!

	// start the snapshot and shutdown the instance (if requested)
	logging.LogMessage(context.IOStream, "Creating snapshot of instance %s", bi.Positionals.Slug)

	// create a new temporary directory
	logging.LogMessage(context.IOStream, "Creating new snapshot staging directory")
	sPath, err := dis.NewSnapshotStagingDir(instance.Slug)
	if err != nil {
		return errSnapshotFailed.Wrap(err)
	}
	defer func() {
		logging.LogMessage(context.IOStream, "Removing snapshot staging directory")
		os.RemoveAll(sPath)
	}()

	// take a snapshot into the staging area!
	sreport := instance.Snapshot(context.IOStream, bi.Keepalive, sPath)

	// write out the report!
	logging.LogOperation(func() error {

		logging.LogOperation(func() error {
			reportPath := filepath.Join(sPath, "report.txt")
			context.Println(reportPath)

			// create the report file!
			report, err := os.Create(reportPath)
			if err != nil {
				return err
			}
			defer report.Close()

			// print the report into it!
			_, err = fmt.Fprintf(report, "%#v\n", sreport)
			return err
		}, context.IOStream, "Writing snapshot report")

		return nil
	}, context.IOStream, "Creating snapshot")

	// copy everything into the final archive

	finalPath := bi.Positionals.Dest
	if finalPath == "" {
		finalPath = dis.NewSnapshotArchivePath(instance.Slug)
	}

	if err := logging.LogOperation(func() error {
		context.IOStream.Println(finalPath)

		targz.Package(finalPath, sPath, func(src string) {
			context.Println(src)
		})
		return err
	}, context.IOStream, "Writing final backup"); err != nil {
		return errSnapshotFailed.Wrap(err)
	}
	context.Printf("Wrote %s\n", finalPath)

	return nil
}
