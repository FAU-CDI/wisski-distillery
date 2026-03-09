package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"go.tkw01536.de/pkglib/errorsx"
)

type TriplestoreUserPayload struct {
	Password           string                     `json:"password"`
	AppSettings        TriplestoreUserAppSettings `json:"appSettings"`
	GrantedAuthorities []string                   `json:"grantedAuthorities"`
}
type TriplestoreUserAppSettings struct {
	DefaultInference      bool `json:"DEFAULT_INFERENCE"`
	DefaultVisGraphSchema bool `json:"DEFAULT_VIS_GRAPH_SCHEMA"`
	DefaultSameas         bool `json:"DEFAULT_SAMEAS"`
	IgnoreSharedQueries   bool `json:"IGNORE_SHARED_QUERIES"`
	ExecuteCount          bool `json:"EXECUTE_COUNT"`
}

var ErrUpdateUserUnauthorized = errors.New("access denied")
var errCreateUserWrongStatusCode = errors.New("failed to create triplestore user: endpoint request did not return status code 201 Created")

// CreateUser creates a new user with the given username and payload.
func (client *Client) CreateUser(ctx context.Context, user string, update TriplestoreUserPayload) (e error) {
	res, err := client.doRestWithMarshal(ctx, http.MethodPost, "/rest/security/users/"+url.PathEscape(user), nil, update)
	if err != nil {
		return fmt.Errorf("failed to create triplestore user: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create triplestore user: %w", errCreateWrongStatusCode)
	}
	return nil
}

// UpdateUser updates the given username to have the provided password.
// If access is denied, return ErrUpdateUserUnauthorized.
func (client *Client) UpdateUser(ctx context.Context, user string, update TriplestoreUserPayload) (e error) {
	res, err := client.doRestWithMarshal(ctx, http.MethodPut, "/rest/security/users/"+url.PathEscape(user), nil, update)
	if err != nil {
		return fmt.Errorf("failed to create triplestore user: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	switch {
	case res.StatusCode == http.StatusOK:
		// we set the password => requests are unauthorized
		// so we still need to enable security (see below!)
		return nil
	case res.StatusCode == http.StatusUnauthorized:
		// a password is needed => security is already enabled.
		// the password may or may not work, but that's a problem for later
		return ErrUpdateUserUnauthorized
	default:
		return fmt.Errorf("failed to create triplestore user: %w", err)
	}
}
