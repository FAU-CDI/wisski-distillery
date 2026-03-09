package client

import (
	"context"
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

// CreateUser creates a new user with the given username and payload.
func (client *Client) CreateUser(ctx context.Context, user string, update TriplestoreUserPayload) (e error) {
	res, err := client.doRestWithMarshal(ctx, http.MethodPost, "/rest/security/users/"+url.PathEscape(user), headers{}, update)
	if err != nil {
		return fmt.Errorf("failed to send http request to users endpoint: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if err := newStatusError(res, true, http.StatusCreated); err != nil {
		return fmt.Errorf("users endpoint responded: %w", err)
	}
	return nil
}

// UpdateUser updates the given username to have the provided password.
func (client *Client) UpdateUser(ctx context.Context, user string, update TriplestoreUserPayload) (e error) {
	res, err := client.doRestWithMarshal(ctx, http.MethodPut, "/rest/security/users/"+url.PathEscape(user), headers{}, update)
	if err != nil {
		return fmt.Errorf("failed to send http request to users endpoint: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if err := newStatusError(res, true, http.StatusOK); err != nil {
		return fmt.Errorf("users endpoint responded: %w", err)
	}
	return nil
}
