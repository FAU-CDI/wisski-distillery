package logging

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/term"
)

func Log[T any](operation func() T, name string, context context.Context) (res T) {
	var took time.Duration

	logger := zerolog.Ctx(context)
	logger.Log().Msg(name)
	defer func() {
		logger.Log().Dur("took", took).Msg(name)
	}()

	start := time.Now()
	res = operation()
	took = time.Since(start)

	return
}

// LogOperation logs a message that is displayed to the user, and then increases the log indent level.
func LogOperation(operation func() error, progress io.Writer, ctx context.Context, format string, args ...interface{}) error {
	logOperation(progress, ctx, getIndent(progress), format, args...)
	incIndent(progress)
	defer decIndent(progress)

	return operation()
}

// Progress writes a progress message to the given progress writer.
func Progress(progress io.Writer, ctx context.Context, message string) {
	io.WriteString(progress, message)
}

// ProgressF is like progress, but uses fmt.Sprintf()
func ProgressF(progress io.Writer, ctx context.Context, format string, args ...interface{}) {
	Progress(progress, ctx, fmt.Sprintf(format, args...))
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
