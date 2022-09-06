package env

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/FAU-CDI/wisski-distillery/internal/fsx"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/tkw1536/goprogram/stream"
)

type SnapshotReport struct {
	Keepalive bool        // was the instance alive while running the snapshot?
	Panic     interface{} // was there a panic?

	// errors for the various components of the Snapshot
	StopErr        error
	StartErr       error
	BookkeepingErr error
	FilesystemErr  error
	TriplestoreErr error
	SQLErr         error
}

// Snapshot creates a new snapshot of this instance into dest
func (instance Instance) Snapshot(io stream.IOStream, keepalive bool, dest string) (report SnapshotReport) {
	// catch anything critical that happened during the snapshot
	defer func() {
		report.Panic = recover()
	}()

	stack := instance.Stack()

	// stop the instance (unless it was explicitly asked to not do so!)
	report.Keepalive = keepalive
	if !keepalive {
		logging.LogMessage(io, "Stopping instance")
		report.StopErr = stack.Down(io)
		defer func() {
			logging.LogMessage(io, "Starting instance")
			report.StartErr = stack.Up(io)
		}()
	}

	// create a wait group, and message channel
	wg := &sync.WaitGroup{}
	messages := make(chan string, 4)

	// write bookkeeping information
	wg.Add(1)
	go func() {
		defer wg.Done()

		bkPath := filepath.Join(dest, "bookkeeping.txt")
		messages <- bkPath

		info, err := os.Create(bkPath)
		if err != nil {
			report.BookkeepingErr = err
			return
		}
		defer info.Close()

		// print whatever is in the database
		// TODO: This should be sql code, maybe gorm can do that?
		_, report.BookkeepingErr = fmt.Fprintf(info, "%#v\n", instance.Instance)
	}()

	// backup the filesystem
	wg.Add(1)
	go func() {
		defer wg.Done()

		fsPath := filepath.Join(dest, filepath.Base(instance.FilesystemBase))
		if err := os.Mkdir(fsPath, fs.ModeDir); err != nil {
			report.FilesystemErr = err
			return
		}

		// copy over whatever is in the base directory
		report.FilesystemErr = fsx.CopyDirectory(fsPath, instance.FilesystemBase, func(dst, src string) {
			messages <- dst
		})
	}()

	// backup the graph db repository
	wg.Add(1)
	go func() {
		defer wg.Done()

		tsPath := filepath.Join(dest, instance.GraphDBRepository+".nq")
		messages <- tsPath

		nquads, err := os.Create(tsPath)
		if err != nil {
			report.TriplestoreErr = err
		}
		defer nquads.Close()

		// directly store the result
		_, report.TriplestoreErr = instance.dis.Triplestore().Backup(nquads, instance.GraphDBRepository)
	}()

	// backup the sql database
	wg.Add(1)
	go func() {
		defer wg.Done()

		sqlPath := filepath.Join(dest, instance.SqlDatabase+".sql")
		messages <- sqlPath

		sql, err := os.Create(sqlPath)
		if err != nil {
			report.SQLErr = err
			return
		}
		defer sql.Close()

		// directly store the result
		report.SQLErr = instance.dis.SQL().Backup(io, sql, instance.SqlDatabase)
	}()

	// TODO: Backup the docker image

	// wait for the group, then close the message channel.
	go func() {
		wg.Wait()
		close(messages)
	}()

	// print out all the messages
	for message := range messages {
		io.Println(message)
	}

	return
}
