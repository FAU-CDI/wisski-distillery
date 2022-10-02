package snapshots

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/pkg/countwriter"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/stream"
)

func (snapshot Snapshot) String() string {
	var builder strings.Builder
	snapshot.Report(&builder)
	return builder.String()
}

// Report writes a report from snapshot into w
func (snapshot Snapshot) Report(w io.Writer) (int, error) {
	ww := countwriter.NewCountWriter(w)

	encoder := json.NewEncoder(ww)
	encoder.SetIndent("", "  ")

	io.WriteString(ww, "======= Begin Snapshot Report "+snapshot.Instance.Slug+" =======\n")

	fmt.Fprintf(ww, "Slug:  %s\n", snapshot.Instance.Slug)
	fmt.Fprintf(ww, "Dest:  %s\n", snapshot.Description.Dest)

	fmt.Fprintf(ww, "Start: %s\n", snapshot.StartTime)
	fmt.Fprintf(ww, "End:   %s\n", snapshot.EndTime)
	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= Description =======\n")
	encoder.Encode(snapshot.Description)
	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= Instance =======\n")
	encoder.Encode(snapshot.Instance)
	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= Errors =======\n")
	fmt.Fprintf(ww, "Panic:       %v\n", snapshot.ErrPanic)
	fmt.Fprintf(ww, "Start:       %s\n", snapshot.ErrStart)
	fmt.Fprintf(ww, "Stop:        %s\n", snapshot.ErrStop)

	fmt.Fprintf(ww, "Whitebox:    %s\n", snapshot.ErrWhitebox)
	fmt.Fprintf(ww, "Blackbox:    %s\n", snapshot.ErrBlackbox)

	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= Manifest =======\n")
	for _, file := range snapshot.Manifest {
		io.WriteString(ww, file+"\n")
	}

	io.WriteString(ww, "\n")

	io.WriteString(ww, "======= End Snapshot Report "+snapshot.Instance.Slug+"=======\n")

	return ww.Sum()
}

// WriteReport writes out the report belonging to this snapshot.
// It is a separate function, to allow writing it indepenently of the rest.
func (snapshot *Snapshot) WriteReport(env environment.Environment, stream stream.IOStream) error {
	return logging.LogOperation(func() error {
		reportPath := filepath.Join(snapshot.Description.Dest, "report.txt")
		stream.Println(reportPath)

		// create the report file!
		report, err := env.Create(reportPath, environment.DefaultFilePerm)
		if err != nil {
			return err
		}
		defer report.Close()

		// print the report into it!
		_, err = io.WriteString(report, snapshot.String())
		return err
	}, stream, "Writing snapshot report")
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
