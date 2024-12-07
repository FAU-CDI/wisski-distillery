//spellchecker:words exporter
package exporter

//spellchecker:words encoding json strings github pkglib sequence
import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/tkw1536/pkglib/sequence"
)

func (snapshot Snapshot) String() string {
	var builder strings.Builder

	_ = snapshot.ReportPlain(&builder) // no way to report error
	return builder.String()
}

func (snapshot Snapshot) ReportMachine(w io.Writer) error {
	return json.NewEncoder(w).Encode(snapshot)
}

//nolint:errcheck
func (snapshot Snapshot) ReportPlain(w io.Writer) error {
	ww := &sequence.Writer{Writer: w} // allows us to ignore all the errors

	encoder := json.NewEncoder(ww)
	encoder.SetIndent("", "  ")

	_, _ = io.WriteString(ww, "======= Begin Snapshot Report "+snapshot.Instance.Slug+" =======\n")

	_, _ = fmt.Fprintf(ww, "Slug:  %s\n", snapshot.Instance.Slug)
	_, _ = fmt.Fprintf(ww, "Dest:  %s\n", snapshot.Description.Dest)

	_, _ = fmt.Fprintf(ww, "Start: %s\n", snapshot.StartTime)
	_, _ = fmt.Fprintf(ww, "End:   %s\n", snapshot.EndTime)
	_, _ = io.WriteString(ww, "\n")

	_, _ = io.WriteString(ww, "======= Description =======\n")
	_ = encoder.Encode(snapshot.Description)
	_, _ = io.WriteString(ww, "\n")

	_, _ = io.WriteString(ww, "======= Instance =======\n")
	_ = encoder.Encode(snapshot.Instance)
	_, _ = io.WriteString(ww, "\n")

	_, _ = io.WriteString(ww, "======= Errors =======\n")
	_, _ = fmt.Fprintf(ww, "Panic:       %v\n", snapshot.ErrPanic)
	_, _ = fmt.Fprintf(ww, "Start:       %v\n", snapshot.ErrStart)
	_, _ = fmt.Fprintf(ww, "Stop:        %v\n", snapshot.ErrStop)

	_, _ = fmt.Fprintf(ww, "Errors:    %s\n", snapshot.Errors)

	_, _ = io.WriteString(ww, "\n")

	_, _ = io.WriteString(ww, "======= Manifest =======\n")
	for _, file := range snapshot.Manifest {
		_, _ = io.WriteString(ww, file+"\n")
	}

	_, _ = io.WriteString(ww, "\n")

	_, _ = io.WriteString(ww, "======= End Snapshot Report "+snapshot.Instance.Slug+"=======\n")

	_, err := ww.Sum()
	return err
}

// Strings turns this backup into a string for the BackupReport.
func (backup Backup) String() string {
	var builder strings.Builder
	_ = backup.ReportPlain(&builder) // no way to report error
	return builder.String()
}

func (backup Backup) ReportMachine(w io.Writer) error {
	return json.NewEncoder(w).Encode(backup)
}

// Report formats a report for this backup, and writes it into Writer.
//
//nolint:errcheck
func (backup Backup) ReportPlain(w io.Writer) error {
	cw := &sequence.Writer{Writer: w}

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

	_, err := cw.Sum()
	return err
}
