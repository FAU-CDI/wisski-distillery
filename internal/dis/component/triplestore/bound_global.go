package triplestore

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore/client"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"go.tkw01536.de/pkglib/errorsx"
)

// boundGlobal implements a wrapper around the global triplestore client.
type boundGlobal struct {
	client   *client.Client
	instance models.Instance
}

func (bound *boundGlobal) ReadURL() string {
	return "http://triplestore:7200/repositories/" + url.PathEscape(bound.instance.GraphDBRepository)
}

func (bound *boundGlobal) WriteURL() string {
	return "http://triplestore:7200/repositories/" + url.PathEscape(bound.instance.GraphDBRepository) + "/statements"
}

func (bound *boundGlobal) Credentials() (username string, password string) {
	return bound.instance.GraphDBUsername, bound.instance.GraphDBPassword
}

// RestoreDB snapshots the provided repository into dst.
func (bound *boundGlobal) RestoreDB(ctx context.Context, progress io.Writer, reader io.Reader) (e error) {
	if err := bound.client.ReplaceContent(ctx, bound.instance.GraphDBRepository, reader); err != nil {
		return fmt.Errorf("failed to restore content: %w", err)
	}
	return nil
}

// Purge purges the given repository and user.
func (bound *boundGlobal) Purge(ctx context.Context, progress io.Writer, allowCreate bool) error {
	return errorsx.Combine(
		bound.client.DeleteRepository(ctx, bound.instance.GraphDBRepository),
		bound.client.DeleteUser(ctx, bound.instance.GraphDBUsername),
	)
}

// SnapshotDB snapshots the provided repository into dst.
func (bound *boundGlobal) SnapshotDB(ctx context.Context, progress io.Writer, dst io.Writer) error {
	_, err := bound.client.ExportContent(ctx, dst, bound.instance.GraphDBRepository)
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to export content: %w", err)
}

// Provision provisions the repository for this instance, possibly deleting any existing repositories.
func (bound *boundGlobal) Provision(ctx context.Context, progress io.Writer, domain string) (e error) {
	if err := bound.client.Wait(ctx, progress); err != nil {
		return fmt.Errorf("failed to wait for triplestore to be ready: %w", err)
	}

	// create the repository
	if err := bound.client.CreateRepository(ctx, client.CreateOpts{
		RepositoryID: bound.instance.GraphDBRepository,
		Label:        domain,
		BaseURL:      "http://" + domain + "/",
	}); err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// create the user and grant them access
	if err := bound.client.CreateUser(ctx, bound.instance.GraphDBUsername, client.TriplestoreUserPayload{
		Password: bound.instance.GraphDBPassword,
		AppSettings: client.TriplestoreUserAppSettings{
			DefaultInference:      true,
			DefaultVisGraphSchema: true,
			DefaultSameas:         true,
			IgnoreSharedQueries:   false,
			ExecuteCount:          true,
		},
		GrantedAuthorities: []string{
			"ROLE_USER",
			"READ_REPO_" + bound.instance.GraphDBRepository,
			"WRITE_REPO_" + bound.instance.GraphDBRepository,
		},
	}); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
