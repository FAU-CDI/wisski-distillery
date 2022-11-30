package triplestore

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/pkg/errors"
)

var errTriplestoreFailedSecurity = errors.New("failed to enable triplestore security: request did not succeed with HTTP 200 OK")

func (ts Triplestore) Update(ctx context.Context, progress io.Writer) error {
	logging.LogMessage(progress, "Waiting for Triplestore")
	if err := ts.Wait(ctx); err != nil {
		return err
	}

	logging.LogMessage(progress, "Resetting admin user password")
	{
		res, err := ts.OpenRaw(ctx, "PUT", "/rest/security/users/"+ts.Config.TriplestoreAdminUser, TriplestoreUserPayload{
			Password: ts.Config.TriplestoreAdminPassword,
			AppSettings: TriplestoreUserAppSettings{
				DefaultInference:      true,
				DefaultVisGraphSchema: true,
				DefaultSameas:         true,
				IgnoreSharedQueries:   false,
				ExecuteCount:          true,
			},
			GrantedAuthorities: []string{"ROLE_ADMIN"},
		}, "", "")
		if err != nil {
			return fmt.Errorf("failed to create triplestore user: %s", err)
		}
		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusOK:
			// we set the password => requests are unauthorized
			// so we still need to enable security (see below!)
		case http.StatusUnauthorized:
			// a password is needed => security is already enabled.
			// the password may or may not work, but that's a problem for later
			logging.LogMessage(progress, "Security is already enabled")
			return nil
		default:
			return fmt.Errorf("failed to create triplestore user: %s", err)
		}
	}

	logging.LogMessage(progress, "Enabling Triplestore security")
	{
		res, err := ts.OpenRaw(ctx, "POST", "/rest/security", true, "", "")
		if err != nil {
			return fmt.Errorf("failed to enable triplestore security: %s", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return errTriplestoreFailedSecurity
		}

		return nil
	}
}
