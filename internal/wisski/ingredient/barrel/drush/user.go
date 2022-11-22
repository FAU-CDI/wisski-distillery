package drush

import (
	"strings"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

var errLoginFailed = exit.Error{
	Message:  "Failed to login",
	ExitCode: exit.ExitGeneric,
}

// Login generates a one-time login url for the given user
func (drush *Drush) Login(io stream.IOStream, user string) (string, error) {
	var builder strings.Builder

	url := drush.Liquid.URL().String()
	command := shellescape.QuoteCommand([]string{"drush", "user:login", "--name=" + user, "--no-browser", "--uri=" + url})

	code, err := drush.Barrel.Shell(io.Streams(&builder, nil, nil, 0), "-c", command)
	if code != 0 || err != nil {
		return "", errLoginFailed
	}
	return strings.TrimSpace(builder.String()), nil
}

var errSetPasswordFailed = exit.Error{
	Message:  "Failed to set password",
	ExitCode: exit.ExitGeneric,
}

func (drush *Drush) ResetPassword(io stream.IOStream, user, password string) error {
	code, err := drush.Barrel.Shell(io, "-c", "drush", "user:password", user, password)
	if code != 0 || err != nil {
		return errSetPasswordFailed
	}
	return nil
}
