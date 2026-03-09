package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"go.tkw01536.de/pkglib/errorsx"
)

// SetSecurity enables or disables the security option of the triplestore.
func (client *Client) SetSecurity(ctx context.Context, enabled bool) (e error) {
	res, err := client.doRestWithMarshal(ctx, http.MethodPost, "/rest/security", headers{}, enabled)
	if err != nil {
		return fmt.Errorf("failed to send http request to security endpoint: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if err := newStatusError(res, true, http.StatusOK); err != nil {
		return fmt.Errorf("security endpoint responded: %w", err)
	}

	return nil
}

// purgeUser deletes the specified user from the triplestore.
// When the user does not exist, returns no error.
func (client *Client) DeleteUser(ctx context.Context, user string) (e error) {
	res, err := client.rest(ctx, http.MethodDelete, "/rest/security/users/"+url.PathEscape(user), headers{})
	if err != nil {
		return fmt.Errorf("failed to send http request to users endpoint: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if err := newStatusError(res, true, http.StatusNoContent, http.StatusNotFound); err != nil {
		return fmt.Errorf("users endpoint responded: %w", err)
	}
	return nil
}
