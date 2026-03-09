//spellchecker:words triplestore
package triplestore

//spellchecker:words bytes context http text template embed github wisski distillery internal models pkglib errorsx
import (
	"context"
	"fmt"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

var errWrongEndpointStatusCode = fmt.Errorf("endpoint request did not return status code %d", http.StatusCreated)

type createRepoContext struct {
	RepositoryID string
	Label        string
	BaseURL      string
}

func (ts *Triplestore) Provision(ctx context.Context, instance models.Instance, domain string, stack *component.StackWithResources) error {
	return ts.For(instance).Provision(ctx, domain)
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
