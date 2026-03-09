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

type createRepoContext struct {
	RepositoryID string
	Label        string
	BaseURL      string
}

func (ts *Client) CreateRepository(ctx context.Context, id, domain, user, password string) (e error) {
	if err := ts.Wait(ctx); err != nil {
		return err
	}

	// prepare the create repo request
	var createRepo bytes.Buffer
	if err := createRepoTemplate.Execute(&createRepo, createRepoContext{
		RepositoryID: id,
		Label:        domain,
		BaseURL:      "http://" + domain + "/",
	}); err != nil {
		return fmt.Errorf("failed to create repository with template: %w", err)
	}

	// do the create!
	{
		res, err := ts.doRestWithForm(ctx, http.MethodPost, "/rest/repositories", nil, "config", &createRepo)
		if err != nil {
			return fmt.Errorf("repository create endpoint failed: %w", err)
		}
		defer errorsx.Close(res.Body, &e, "response body")
		if res.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create repository: %w", errCreateWrongStatusCode)
		}

		return nil
	}
}

// DeleteRepository deletes the specified repo from the triplestore.
// When the repo does not exist, returns no error.
func (client *Client) DeleteRepository(ctx context.Context, repo string) (e error) {
	res, err := client.rest(ctx, http.MethodDelete, "/rest/repositories/"+url.PathEscape(repo), nil)
	if err != nil {
		return err
	}
	defer errorsx.Close(res.Body, &e, "response body")
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("%w: %d", errDeleteUserStatusCode, res.StatusCode)
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
	res, err := client.rest(ctx, http.MethodGet, "/rest/repositories", &requestHeaders{Accept: "application/json"})
	if err != nil {
		return nil, err
	}
	defer errorsx.Close(res.Body, &e, "response body")

	e = json.NewDecoder(res.Body).Decode(&repos)
	return
}
