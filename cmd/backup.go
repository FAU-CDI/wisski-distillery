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

// BackupInstance is the 'backup_instance' command
var BackupInstance wisski_distillery.Command = backupInstance{}

type backupInstance struct {
	Keepalive   bool `short:"k" long:"keepalive" description:"Keep instance running while taking a backup. Might lead to inconsistent state"`
	Positionals struct {
		Slug    string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to show info about"`
		Outfile string `positional-arg-name:"OUTFILE" description:"Destination file to write backup to"`
	} `positional-args:"true"`
}

func (backupInstance) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: env.Requirements{
			NeedsConfig: true,
		},
		Command:     "backup_instance",
		Description: "Makes a backup of a specific instance",
	}
}

var errBackupFailed = exit.Error{
	Message:  "Failed to make a backup",
	ExitCode: exit.ExitGeneric,
}

func (bi backupInstance) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	instance, err := dis.Instance(bi.Positionals.Slug)
	if err != nil {
		return err
	}

	// TODO: Allow skipping backups of individual parts and make them concurrent!

	// start the backup and shutdown the instance (if requested)
	logging.LogMessage(context.IOStream, "Creating backup of instance %s", bi.Positionals.Slug)

	// create a new temporary directory
	logging.LogMessage(context.IOStream, "Creating temporary backup directory")
	path, err := dis.NewInprogressBackupPath(instance.Slug)
	if err != nil {
		return errBackupFailed.Wrap(err)
	}
	defer func() {
		logging.LogMessage(context.IOStream, "Removing temporary backup directory")
		os.RemoveAll(path) // TODO: Turn this on again
	}()

	// make a snapshot and write out the report also!
	logging.LogOperation(func() error {
		sreport := instance.Snapshot(context.IOStream, bi.Keepalive, path)

		logging.LogOperation(func() error {
			reportPath := filepath.Join(path, "report.txt")
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

	// copy everything into the final file!
	finalPath := bi.Positionals.Outfile
	if finalPath == "" {
		finalPath = dis.FinalBackupArchive(instance.Slug)
	}

	if err := logging.LogOperation(func() error {
		context.IOStream.Println(finalPath)

		targz.Package(finalPath, path, func(src string) {
			context.Println(src)
		})
		return err
	}, context.IOStream, "Writing final backup"); err != nil {
		return errBackupFailed.Wrap(err)
	}
	context.Printf("Wrote %s\n", finalPath)

	return nil
}
