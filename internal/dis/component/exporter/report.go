package exporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tkw1536/pkglib/pools"
	"github.com/tkw1536/pkglib/sequence"
)

func (snapshot Snapshot) String() string {
	builder := pools.GetBuilder()
	defer pools.ReleaseBuilder(builder)

	snapshot.Report(builder)
	return builder.String()
}

// Report writes a report from snapshot into w
func (snapshot Snapshot) Report(w io.Writer) (int, error) {
	ww := &sequence.Writer{Writer: w}

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

// Strings turns this backup into a string for the BackupReport.
func (backup Backup) String() string {
	builder := pools.GetBuilder()
	defer pools.ReleaseBuilder(builder)

	backup.Report(builder)
	return builder.String()
}

// Report formats a report for this backup, and writes it into Writer.
func (backup Backup) Report(w io.Writer) (int, error) {
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

	return cw.Sum()
}
