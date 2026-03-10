//spellchecker:words triplestore
package triplestore

//spellchecker:words context errors http github wisski distillery internal component logging pkglib errorsx
import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore/client"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

func (ts *Triplestore) Update(ctx context.Context, progress io.Writer) (e error) {
	cl := ts.globalClient()
	if _, err := logging.LogMessage(progress, "Waiting for Triplestore"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if err := cl.Wait(ctx, progress); err != nil {
		return fmt.Errorf("failed to wait: %w", err)
	}

	if _, err := logging.LogMessage(progress, "Resetting admin user password"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	{
		config := component.GetStill(ts).Config.TS

		err := cl.UpdateUser(ctx, config.AdminUsername, client.TriplestoreUserPayload{
			Password: config.AdminPassword,
			AppSettings: client.TriplestoreUserAppSettings{
				DefaultInference:      true,
				DefaultVisGraphSchema: true,
				DefaultSameas:         true,
				IgnoreSharedQueries:   false,
				ExecuteCount:          true,
			},
			GrantedAuthorities: []string{"ROLE_ADMIN"},
		})

		var errWrongStatus client.WrongStatusError
		if errors.As(err, &errWrongStatus) && errWrongStatus.Got == http.StatusUnauthorized {
			if _, err := logging.LogMessage(progress, "Security is already enabled"); err != nil {
				return fmt.Errorf("failed to log message: %w", err)
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to create triplestore user: %w", err)
		}

		if _, err := logging.LogMessage(progress, "Enabling Triplestore security"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		if err := cl.SetSecurity(ctx, true); err != nil {
			return fmt.Errorf("failed to enable triplestore security: %w", err)
		}
	}

	return nil
}
