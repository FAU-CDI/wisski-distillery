package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
	"github.com/FAU-CDI/wisski-distillery/internal/fsx"
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

	// make the snapshot!
	// TODO: Ignore errors here, and write them into the snapshot instance
	if err := bi.makeSnapshot(context, path, instance); err != nil {
		return errBackupFailed.WithMessageF(err)
	}

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

// makeSnapshot makes a snapshot of the directory into the given directory!
//
// TODO: Return a SnapshotReport object, and only check what was actually copied
func (bi backupInstance) makeSnapshot(context wisski_distillery.Context, path string, instance env.Instance) error {
	dis := context.Environment
	stack := instance.Stack()

	if !bi.Keepalive {
		logging.LogMessage(context.IOStream, "Stopping instance")
		if err := stack.Down(context.IOStream); err != nil {
			return err
		}
		defer func() {
			logging.LogMessage(context.IOStream, "Starting instance")
			stack.Up(context.IOStream)
		}()
	}

	// backup up bookkeeping info!
	if err := logging.LogOperation(func() error {
		bkPath := filepath.Join(path, "bookkeeping.txt")
		context.IOStream.Println(bkPath)

		// create the backup file!
		info, err := os.Create(bkPath)
		if err != nil {
			return err
		}
		defer info.Close()

		// print whatever is in the bookkeeping instance
		_, err = fmt.Fprintf(info, "%#v\n", instance.Instance)
		return err
	}, context.IOStream, "Backing up Bookkeping Information"); err != nil {
		return errBackupFailed.Wrap(err)
	}

	// backup the filesystem!
	if err := logging.LogOperation(func() error {
		// create a backup directory
		fsPath := filepath.Join(path, filepath.Base(instance.FilesystemBase))
		if err := os.Mkdir(fsPath, fs.ModeDir); err != nil {
			return err
		}

		return fsx.CopyDirectory(fsPath, instance.FilesystemBase, func(dst, src string) {
			context.IOStream.Println(src)
		})
	}, context.IOStream, "Backing up filesystem"); err != nil {
		return errBackupFailed.Wrap(err)
	}

	// backup the the triplestore!
	if err := logging.LogOperation(func() error {
		tsPath := filepath.Join(path, instance.GraphDBRepository+".nq")
		context.IOStream.Println(tsPath)

		// create the backup file!
		nquads, err := os.Create(tsPath)
		if err != nil {
			return err
		}
		defer nquads.Close()

		// TODO: Add a progress bar?
		_, err = dis.Triplestore().Backup(nquads, instance.GraphDBRepository)
		return err
	}, context.IOStream, "Backing up Triplestore"); err != nil {
		return errBackupFailed.Wrap(err)
	}

	// backup the the sql database!
	if err := logging.LogOperation(func() error {
		sqlPath := filepath.Join(path, instance.SqlDatabase+".sql")
		context.IOStream.Println(sqlPath)

		// create the backup file!
		sql, err := os.Create(sqlPath)
		if err != nil {
			return err
		}
		defer sql.Close()

		// TODO: Add a progress bar?
		return dis.SQL().Backup(context.IOStream, sql, instance.SqlDatabase)
	}, context.IOStream, "Backing up Triplestore"); err != nil {
		return errBackupFailed.Wrap(err)
	}

	return nil
}
