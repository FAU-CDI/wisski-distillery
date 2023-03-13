package logging

import (
	"context"
	"fmt"
	"io"
	"strings"

	"golang.org/x/term"
)

// LogOperation logs a message that is displayed to the user, and then increases the log indent level.
func LogOperation(operation func() error, progress io.Writer, ctx context.Context, format string, args ...interface{}) error {
	logOperation(progress, ctx, getIndent(progress), format, args...)
	incIndent(progress)
	defer decIndent(progress)

	return operation()
}

// LogMessage logs a message that is displayed to the user
func LogMessage(progress io.Writer, ctx context.Context, format string, args ...interface{}) (int, error) {
	return logOperation(progress, ctx, getIndent(progress), format, args...)
}

func logOperation(progress io.Writer, ctx context.Context, indent int, format string, args ...interface{}) (int, error) {
	message := "\033[1m" + strings.Repeat(" ", indent+1) + "=> " + format + "\033[0m\n"
	if !streamIsTerminal(progress) {
		message = " => " + format + "\n"
	}

	return fmt.Fprintf(progress, message, args...)
}

// streamIsTerminal checks if stream is a terminal
func streamIsTerminal(stream any) bool {
	file, ok := stream.(interface{ Fd() uintptr })
	return ok && term.IsTerminal(int(file.Fd()))
}
