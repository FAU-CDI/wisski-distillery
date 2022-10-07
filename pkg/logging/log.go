package logging

import (
	"strings"

	"github.com/tkw1536/goprogram/stream"
)

// LogOperation logs a message that is displayed to the user, and then increases the log indent level.
func LogOperation(operation func() error, io stream.IOStream, format string, args ...interface{}) error {
	logOperation(io, getIndent(io), format, args...)
	incIndent(io)
	defer decIndent(io)

	return operation()
}

// LogMessage logs a message that is displayed to the user
func LogMessage(io stream.IOStream, format string, args ...interface{}) (int, error) {
	return logOperation(io, getIndent(io), format, args...)
}

func logOperation(io stream.IOStream, indent int, format string, args ...interface{}) (int, error) {
	message := "\033[1m" + strings.Repeat(" ", indent+1) + "=> " + format + "\033[0m\n"
	if !io.StdoutIsATerminal() {
		message = " => " + format + "\n"
	}

	return io.Printf(message, args...)
}
