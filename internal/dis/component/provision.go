//spellchecker:words component
package component

//spellchecker:words context github wisski distillery internal models
import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// Provisionable is a component that can be provisioned.
type Provisionable interface {
	Component

	// ProvisionNeedsStack indicates if this provisionable should be provisioned after the inital stack is set up.
	ProvisionNeedsStack(instance models.Instance) bool

	// PurgeMayFail indicates if this provisionable may fail to purge.
	PurgeMayFail(instance models.Instance) bool

	// Provision provisions resources specific to the provided instance.
	//
	// Domain holds the full (unique) domain name of the given instance.
	//
	// If stack is nil, it is guaranteed that ProvisionNeedsStack() was called and returned false.
	// If stack is not nil, either ProvisionNeedsStack() was called and returned true, or it was not called at all.
	Provision(ctx context.Context, progress io.Writer, instance models.Instance, domain string, stack *StackWithResources) error

	// Purge purges resources specific to the provided instance.
	// Domain holds the full (unique) domain name of the given instance.
	Purge(ctx context.Context, progress io.Writer, instance models.Instance, domain string) error
}
