package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"go.tkw01536.de/pkglib/timex"
)

// Wait waits for the connection to the Triplestore to succeed.
// This is achieved using a polling strategy.
func (client *Client) Wait(ctx context.Context) error {
	if err := timex.TickUntilFunc(func(time.Time) bool {
		res, err := client.rest(ctx, http.MethodGet, "/rest/repositories", headers{})
		wdlog.Of(ctx).Debug(
			"Triplestore Wait",
			"error", err,
		)
		if err != nil {
			return false
		}

		defer res.Body.Close() //nolint:errcheck // no way to report error
		return true
	}, ctx, client.PollInterval); err != nil {
		return fmt.Errorf("failed to wait for triplestore: %w", err)
	}
	return nil
}
