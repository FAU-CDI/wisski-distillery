//spellchecker:words logging
package logging

//spellchecker:words strings github pkglib errorsx golang term
import (
	"fmt"
	"io"
	"strings"

	"go.tkw01536.de/pkglib/errorsx"
	"golang.org/x/term"
)

// LogOperation logs a message that is displayed to the user, and then increases the log indent level.
func LogOperation(operation func() error, progress io.Writer, format string, args ...interface{}) error {
	_, errLog := logOperation(progress, getIndent(progress), format, args...)
	incIndent(progress)
	defer decIndent(progress)

	err := operation()
	return errorsx.Combine(err, errLog)
}

// LogMessage logs a message that is displayed to the user.
func LogMessage(progress io.Writer, format string, args ...interface{}) (int, error) {
	return logOperation(progress, getIndent(progress), format, args...)
}

func logOperation(progress io.Writer, indent int, format string, args ...interface{}) (int, error) {
	message := "\033[1m" + strings.Repeat(" ", indent+1) + "=> " + format + "\033[0m\n"
	if !streamIsTerminal(progress) {
		message = " => " + format + "\n"
	}

	count, err := fmt.Fprintf(progress, message, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to format message: %w", err)
	}
	return count, nil
}

// streamIsTerminal checks if stream is a terminal.
func streamIsTerminal(stream any) bool {
	file, ok := stream.(interface{ Fd() uintptr })
	return ok && term.IsTerminal(int(file.Fd()))
}
