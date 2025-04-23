//spellchecker:words triplestore
package triplestore

//spellchecker:words bytes context errors http text template embed github wisski distillery internal models goprogram exit
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/goprogram/exit"
)

var errTripleStoreFailedRepository = exit.Error{
	Message:  "failed to create repository: %s",
	ExitCode: exit.ExitGeneric,
}

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

func (ts *Triplestore) Provision(ctx context.Context, instance models.Instance, domain string) error {
	return ts.CreateRepository(ctx, instance.GraphDBRepository, domain, instance.GraphDBUsername, instance.GraphDBPassword)
}

func (ts *Triplestore) Purge(ctx context.Context, instance models.Instance, domain string) error {
	return errors.Join(
		ts.PurgeRepo(ctx, instance.GraphDBRepository),
		ts.PurgeUser(ctx, instance.GraphDBUsername),
	)
}

func (ts *Triplestore) CreateRepository(ctx context.Context, name, domain, user, password string) error {
	if err := ts.Wait(ctx); err != nil {
		return err
	}

	// prepare the create repo request
	var createRepo bytes.Buffer
	if err := createRepoTemplate.Execute(&createRepo, createRepoContext{
		RepositoryID: name,
		Label:        domain,
		BaseURL:      "http://" + domain + "/",
	}); err != nil {
		return fmt.Errorf("failed to create repository with template: %w", err)
	}

	// do the create!
	{
		res, err := ts.DoRestWithForm(ctx, tsTrivialTimeout, http.MethodPost, "/rest/repositories", nil, "config", &createRepo)
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
		res, err := ts.DoRestWithMarshal(ctx, tsTrivialTimeout, http.MethodPost, "/rest/security/users/"+user, nil, TriplestoreUserPayload{
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
		})
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
