package triplestore

import (
	"bytes"
	"context"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/errorx"
	"github.com/tkw1536/pkglib/pools"
)

var errTripleStoreFailedRepository = exit.Error{
	Message:  "failed to create repository: %s",
	ExitCode: exit.ExitGeneric,
}

//go:embed create-repo.ttl
var createRepoTTL []byte

func (ts *Triplestore) Provision(ctx context.Context, instance models.Instance, domain string) error {
	return ts.CreateRepository(ctx, instance.GraphDBRepository, domain, instance.GraphDBUsername, instance.GraphDBPassword)
}

func (ts *Triplestore) Purge(ctx context.Context, instance models.Instance, domain string) error {
	return errorx.First(
		ts.PurgeRepo(ctx, instance.GraphDBRepository),
		ts.PurgeUser(ctx, instance.GraphDBUsername),
	)
}

func (ts *Triplestore) CreateRepository(ctx context.Context, name, domain, user, password string) error {
	if err := ts.Wait(ctx); err != nil {
		return err
	}

	// prepare the create repo request
	createRepo := pools.GetBuffer()
	defer pools.ReleaseBuffer(createRepo)
	err := unpack.WriteTemplate(createRepo, map[string]string{
		"GRAPHDB_REPO":    name,
		"INSTANCE_DOMAIN": domain,
	}, bytes.NewReader(createRepoTTL))
	if err != nil {
		return err
	}

	// do the create!
	{
		res, err := ts.OpenRaw(ctx, "POST", "/rest/repositories", createRepo.Bytes(), "config", "")
		if err != nil {
			return errTripleStoreFailedRepository.WithMessageF(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			return errTripleStoreFailedRepository.WithMessageF("repo create did not return status code 201")
		}
	}

	// create the user and grant them access
	{
		res, err := ts.OpenRaw(ctx, "POST", "/rest/security/users/"+user, TriplestoreUserPayload{
			Password: password,
			AppSettings: TriplestoreUserAppSettings{
				DefaultInference:      true,
				DefaultVisGraphSchema: true,
				DefaultSameas:         true,
				IgnoreSharedQueries:   false,
				ExecuteCount:          true,
			},
			GrantedAuthorities: []string{
				"ROLE_USER",
				"READ_REPO_" + name,
				"WRITE_REPO_" + name,
			},
		}, "", "")
		if err != nil {
			return errTripleStoreFailedRepository.WithMessageF(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			return errTripleStoreFailedRepository.WithMessageF("user create did not return status code 201")
		}
	}

	return nil
}
