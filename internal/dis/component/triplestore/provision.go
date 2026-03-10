//spellchecker:words triplestore
package triplestore

//spellchecker:words bytes context http text template embed github wisski distillery internal models pkglib errorsx
import (
	"context"
	"fmt"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"go.tkw01536.de/pkglib/stream"
)

func (ts *Triplestore) Provision(ctx context.Context, instance models.Instance, domain string, stack *component.StackWithResources) error {
	if err := ts.For(instance).Provision(ctx, stream.Null, domain); err != nil {
		return fmt.Errorf("failed to provision triplestore: %w", err)
	}
	return nil
}

func (ts *Triplestore) ProvisionNeedsStack(instance models.Instance) bool {
	return instance.DedicatedTriplestore
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
