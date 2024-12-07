//spellchecker:words users
package users

//spellchecker:words github wisski distillery internal status ingredient
import (
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

type UserPolicy struct {
	ingredient.Base
}

var (
	_ ingredient.WissKIFetcher = (*UserPolicy)(nil)
)

func (up *UserPolicy) Fetch(flags ingredient.FetcherFlags, target *status.WissKI) (err error) {
	if flags.Quick {
		return nil
	}

	// read the grants into the info struct
	liquid := ingredient.GetLiquid(up)
	target.Grants, err = liquid.Policy.Instance(flags.Context, liquid.Slug)
	return err
}
