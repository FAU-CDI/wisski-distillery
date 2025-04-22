//spellchecker:words triplestore
package triplestore

//spellchecker:words context http github wisski distillery internal component logging errors
import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/pkg/errors"
)

var errTriplestoreFailedSecurity = errors.New("failed to enable triplestore security: request did not succeed with HTTP 200 OK")

func (ts *Triplestore) Update(ctx context.Context, progress io.Writer) error {
	logging.LogMessage(progress, "Waiting for Triplestore")
	if err := ts.Wait(ctx); err != nil {
		return err
	}

	logging.LogMessage(progress, "Resetting admin user password")
	{
		config := component.GetStill(ts).Config.TS

		res, err := ts.DoRestWithMarshal(ctx, tsTrivialTimeout, http.MethodPut, "/rest/security/users/"+config.AdminUsername, nil, TriplestoreUserPayload{
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
			return fmt.Errorf("failed to create triplestore user: %w", err)
		}
	}

	logging.LogMessage(progress, "Enabling Triplestore security")
	{
		res, err := ts.DoRestWithMarshal(ctx, tsTrivialTimeout, http.MethodPost, "/rest/security", nil, true)
		if err != nil {
			return fmt.Errorf("failed to enable triplestore security: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return errTriplestoreFailedSecurity
		}

		return nil
	}
}
