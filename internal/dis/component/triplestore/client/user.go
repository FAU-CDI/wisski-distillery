package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"go.tkw01536.de/pkglib/errorsx"
)

var (
	errTriplestoreFailedSecurity = errors.New("failed to enable triplestore security: request did not succeed with HTTP 200 OK")
	errDeleteUserWrongStatusCode = errors.New("purge returned abnormal exit code")
)

// SetSecurity enables or disables the security option of the triplestore.
func (client *Client) SetSecurity(ctx context.Context, enabled bool) (e error) {
	res, err := client.doRestWithMarshal(ctx, http.MethodPost, "/rest/security", nil, enabled)
	if err != nil {
		return fmt.Errorf("failed to enable triplestore security: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if res.StatusCode != http.StatusOK {
		return errTriplestoreFailedSecurity
	}

	return nil
}

// purgeUser deletes the specified user from the triplestore.
// When the user does not exist, returns no error.
func (client *Client) DeleteUser(ctx context.Context, user string) (e error) {
	res, err := client.rest(ctx, http.MethodDelete, "/rest/security/users/"+url.PathEscape(user), nil)
	if err != nil {
		return err
	}
	defer errorsx.Close(res.Body, &e, "response body")
	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("%w: %d", errDeleteUserWrongStatusCode, res.StatusCode)
	}
	return nil
}
