package ssh

//spellchecker:words context github wisski distillery internal status ingredient gliderlabs golang crypto gossh
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

type SSH struct {
	ingredient.Base
}

var (
	_ ingredient.WissKIFetcher = (*SSH)(nil)
)

func (ssh *SSH) Keys(ctx context.Context) (keys []ssh.PublicKey, err error) {
	liquid := ingredient.GetLiquid(ssh)
	grants, err := liquid.Policy.Instance(ctx, liquid.Slug)
	if err != nil {
		return nil, err
	}

	// iterate over enabled distillery admin users
	for _, grant := range grants {
		if !grant.DrupalAdminRole {
			continue
		}
		ukeys, err := liquid.Keys.Keys(ctx, grant.User)
		if err != nil {
			return nil, err
		}
		for _, ukey := range ukeys {
			if pk := ukey.PublicKey(); pk != nil {
				keys = append(keys, pk)
			}
		}
	}

	// and return the keys!
	return keys, nil
}

// AllKeys returns the keys specifically registered to this instance and all the globally registered keys.
func (ssh *SSH) AllKeys(ctx context.Context) (keys []ssh.PublicKey, err error) {
	lkeys, err := ssh.Keys(ctx)
	if err != nil {
		return nil, err
	}

	gkeys, err := ingredient.GetLiquid(ssh).Keys.Admin(ctx)
	if err != nil {
		return nil, err
	}

	keys = append(keys, lkeys...)
	keys = append(keys, gkeys...)

	return keys, nil
}

func (ssh *SSH) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) error {
	if flags.Quick {
		return nil
	}

	// add the instance keys
	keys, err := ssh.AllKeys(flags.Context)
	if err != nil {
		return err
	}

	info.SSHKeys = make([]string, len(keys))
	for i, key := range keys {
		info.SSHKeys[i] = string(gossh.MarshalAuthorizedKey(key))
	}

	return nil
}
