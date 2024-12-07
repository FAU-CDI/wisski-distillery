//spellchecker:words logging
package logging

//spellchecker:words errors strings golang term
import (
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/term"
)

// LogOperation logs a message that is displayed to the user, and then increases the log indent level.
func LogOperation(operation func() error, progress io.Writer, format string, args ...interface{}) error {
	_, errLog := logOperation(progress, getIndent(progress), format, args...)
	incIndent(progress)
	defer decIndent(progress)

	err := operation()
	if errLog != nil && err != nil {
		return errors.Join(err, errLog)
	}
	return err
}

// LogMessage logs a message that is displayed to the user
func LogMessage(progress io.Writer, format string, args ...interface{}) (int, error) {
	return logOperation(progress, getIndent(progress), format, args...)
}

func logOperation(progress io.Writer, indent int, format string, args ...interface{}) (int, error) {
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
