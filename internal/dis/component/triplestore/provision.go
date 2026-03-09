//spellchecker:words triplestore
package triplestore

//spellchecker:words bytes context http text template embed github wisski distillery internal models pkglib errorsx
import (
	"context"
	"fmt"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore/client"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

var errWrongEndpointStatusCode = fmt.Errorf("endpoint request did not return status code %d", http.StatusCreated)

type createRepoContext struct {
	RepositoryID string
	Label        string
	BaseURL      string
}

func (ts *Triplestore) Provision(ctx context.Context, instance models.Instance, domain string, stack *component.StackWithResources) error {
	return ts.CreateRepository(ctx, instance.GraphDBRepository, domain, instance.GraphDBUsername, instance.GraphDBPassword)
}

func (ts *Triplestore) ProvisionNeedsStack(instance models.Instance) bool {
	return false
}

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

func (ts *Triplestore) CreateRepository(ctx context.Context, name, domain, user, password string) (e error) {
	cl := ts.client()
	if err := cl.Wait(ctx); err != nil {
		return err
	}

	// create the repository
	if err := cl.CreateRepository(ctx, name, domain, user, password); err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// create the user and grant them access
	if err := cl.CreateUser(ctx, user, client.TriplestoreUserPayload{
		Password: password,
		AppSettings: client.TriplestoreUserAppSettings{
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
	}); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
