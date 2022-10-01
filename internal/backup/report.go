package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/pkg/countwriter"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/stream"
)

// Description provides a description for a backup
type Description struct {
	Dest string // Destination path
	Auto bool   // Was the path created automatically?

	ConcurrentSnapshots int // maximum number of concurrent snapshots
}

// Backup describes a backup
type Backup struct {
	Description Description

	// Start and End Time of the backup
	StartTime time.Time
	EndTime   time.Time

	// various error states, which are ignored when creating the snapshot
	ErrPanic interface{}

	// errors for the various components
	ComponentErrors map[string]error

	// TODO: Make this proper
	ConfigFileErr error

	// Snapshots containing instances
	InstanceListErr   error
	InstanceSnapshots []snapshots.Snapshot

	// List of files included
	Manifest []string
}

// Strings turns this backup into a string for the BackupReport.
func (backup Backup) String() string {
	var builder strings.Builder
	backup.Report(&builder)
	return builder.String()
}

// Report formats a report for this backup, and writes it into Writer.
func (backup Backup) Report(w io.Writer) (int, error) {
	cw := countwriter.NewCountWriter(w)

	encoder := json.NewEncoder(cw)
	encoder.SetIndent("", "  ")

	io.WriteString(cw, "======= Backup =======\n")

	fmt.Fprintf(cw, "Start: %s\n", backup.StartTime)
	fmt.Fprintf(cw, "End:   %s\n", backup.EndTime)
	io.WriteString(cw, "\n")

	io.WriteString(cw, "======= Description =======\n")
	encoder.Encode(backup.Description)
	io.WriteString(cw, "\n")

	io.WriteString(cw, "======= Errors =======\n")
	fmt.Fprintf(cw, "Panic:            %v\n", backup.ErrPanic)
	fmt.Fprintf(cw, "Component Errors: %v\n", backup.ComponentErrors)
	fmt.Fprintf(cw, "ConfigFileErr:    %s\n", backup.ConfigFileErr)
	fmt.Fprintf(cw, "InstanceListErr:  %s\n", backup.InstanceListErr)

	io.WriteString(cw, "\n")

	io.WriteString(cw, "======= Snapshots =======\n")
	for _, s := range backup.InstanceSnapshots {
		io.WriteString(cw, s.String())
		io.WriteString(cw, "\n")
	}

	io.WriteString(cw, "======= Manifest =======\n")
	for _, file := range backup.Manifest {
		io.WriteString(cw, file+"\n")
	}

	io.WriteString(cw, "\n")

	return cw.Sum()
}

// WriteReport writes out the report belonging to this backup.
// It is a separate function, to allow writing it indepenently of the rest.
func (backup Backup) WriteReport(env environment.Environment, stream stream.IOStream) error {
	return logging.LogOperation(func() error {
		reportPath := filepath.Join(backup.Description.Dest, "report.txt")
		stream.Println(reportPath)

		// create the report file!
		report, err := env.Create(reportPath, environment.DefaultFilePerm)
		if err != nil {
			return err
		}
		defer report.Close()

		// print the report into it!
		_, err = io.WriteString(report, backup.String())
		return err
	}, stream, "Writing backup report")
}
