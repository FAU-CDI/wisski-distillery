package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"go.tkw01536.de/pkglib/errorsx"

	_ "embed"
)

var (
	errCreateWrongStatusCode = fmt.Errorf("endpoint request did not return status code %d", http.StatusCreated)
	errDeleteUserStatusCode  = errors.New("delete user returned abnormal exit code")
)

//go:embed create-repo.tpl
var createRepoTpl string

// Template for creating repositories.
//
// NOTE(twiesing): The template is not aware of SparQL syntax, thus this template is very unsafe.
// And should only be used with KNOWN GOOD input.
var createRepoTemplate = template.Must(template.New("create-repo.tpl").Parse(createRepoTpl))

type CreateOpts struct {
	RepositoryID string
	Label        string
	BaseURL      string `json:"-"`
}

func (ts *Client) CreateRepository(ctx context.Context, opts CreateOpts) (e error) {
	if err := ts.Wait(ctx); err != nil {
		return fmt.Errorf("failed to wait for repository to be ready: %w", err)
	}

	// prepare the create repo request
	var createRepo bytes.Buffer
	if err := createRepoTemplate.Execute(&createRepo, opts); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// do the create!
	{
		res, err := ts.doRestWithForm(ctx, http.MethodPost, "/rest/repositories", headers{}, "config", &createRepo)
		if err != nil {
			return fmt.Errorf("failed to send http request to repositories endpoint: %w", err)
		}
		defer errorsx.Close(res.Body, &e, "response body")

		if err := newStatusError(res, true, http.StatusCreated); err != nil {
			return fmt.Errorf("repositories endpoint responded: %w", err)
		}
		return nil
	}
}

// DeleteRepository deletes the specified repo from the triplestore.
// When the repo does not exist, returns no error.
func (client *Client) DeleteRepository(ctx context.Context, repo string) (e error) {
	res, err := client.rest(ctx, http.MethodDelete, "/rest/repositories/"+url.PathEscape(repo), headers{})
	if err != nil {
		return fmt.Errorf("failed to send http request to repositories endpoint: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if err := newStatusError(res, true, http.StatusOK, http.StatusNotFound); err != nil {
		return fmt.Errorf("repositories endpoint responded: %w", err)
	}
	return nil
}

type Repository struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	URI        string `json:"uri"`
	Type       string `json:"type"`
	SesameType string `json:"sesameType"`
	Location   string `json:"location"`
	Readable   bool   `json:"readable"`
	Writable   bool   `json:"writable"`
	Local      bool   `json:"local"`
}

func (client *Client) ListRepositories(ctx context.Context) (repos []Repository, e error) {
	res, err := client.rest(ctx, http.MethodGet, "/rest/repositories", headers{Accept: "application/json"})
	if err != nil {
		return nil, fmt.Errorf("failed to send http request to repositories endpoint: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if err := newStatusError(res, true, http.StatusOK); err != nil {
		return nil, fmt.Errorf("repositories endpoint responded: %w", err)
	}

	err = json.NewDecoder(res.Body).Decode(&repos)
	if err != nil {
		return nil, fmt.Errorf("failed to decode repositories: %w", err)
	}
	return repos, nil
}
