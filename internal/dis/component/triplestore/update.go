//spellchecker:words triplestore
package triplestore

//spellchecker:words context errors http github wisski distillery internal component logging pkglib errorsx
import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"go.tkw01536.de/pkglib/errorsx"
)

var errTriplestoreFailedSecurity = errors.New("failed to enable triplestore security: request did not succeed with HTTP 200 OK")

func (ts *Triplestore) Update(ctx context.Context, progress io.Writer) (e error) {
	if _, err := logging.LogMessage(progress, "Waiting for Triplestore"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if err := ts.Wait(ctx); err != nil {
		return fmt.Errorf("failed to wait: %w", err)
	}

	if _, err := logging.LogMessage(progress, "Resetting admin user password"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	{
		config := component.GetStill(ts).Config.TS

		res, err := ts.DoRestWithMarshal(ctx, tsTrivialTimeout, http.MethodPut, "/rest/security/users/"+url.PathEscape(config.AdminUsername), nil, TriplestoreUserPayload{
			Password: config.AdminPassword,
			AppSettings: TriplestoreUserAppSettings{
				DefaultInference:      true,
				DefaultVisGraphSchema: true,
				DefaultSameas:         true,
				IgnoreSharedQueries:   false,
				ExecuteCount:          true,
			},
			GrantedAuthorities: []string{"ROLE_ADMIN"},
		})
		if err != nil {
			return fmt.Errorf("failed to create triplestore user: %w", err)
		}
		defer errorsx.Close(res.Body, &e, "response body")

		switch res.StatusCode {
		case http.StatusOK:
			// we set the password => requests are unauthorized
			// so we still need to enable security (see below!)
		case http.StatusUnauthorized:
			// a password is needed => security is already enabled.
			// the password may or may not work, but that's a problem for later
			if _, err := logging.LogMessage(progress, "Security is already enabled"); err != nil {
				return fmt.Errorf("failed to log message: %w", err)
			}
			return nil
		default:
			return fmt.Errorf("failed to create triplestore user: %w", err)
		}
	}

	if _, err := logging.LogMessage(progress, "Enabling Triplestore security"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	{
		res, err := ts.DoRestWithMarshal(ctx, tsTrivialTimeout, http.MethodPost, "/rest/security", nil, true)
		if err != nil {
			return fmt.Errorf("failed to enable triplestore security: %w", err)
		}
		defer errorsx.Close(res.Body, &e, "response body")

		if res.StatusCode != http.StatusOK {
			return errTriplestoreFailedSecurity
		}

		return nil
	}
}
