package list

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/tkw1536/pkglib/lazy"
	"golang.org/x/sync/errgroup"
)

// ListInstances holds information about all instances
type ListInstances struct {
	component.Base

	names lazy.Lazy[map[string]struct{}] // instance names
	infos lazy.Lazy[[]status.WissKI]     // list of home instances (updated via cron)

	Dependencies struct {
		Auth      *auth.Auth
		Instances *instances.Instances
	}
}

func (li *ListInstances) Names() map[string]struct{} {
	return li.names.Get(nil)
}

func (li *ListInstances) Infos() []status.WissKI {
	return li.infos.Get(nil)
}

// ShouldShowList determines if a list should be shown for the given request
func (li *ListInstances) ShouldShowList(r *http.Request) bool {
	allowPrivate := li.Config.Home.List.Private.Value
	allowPublic := li.Config.Home.List.Public.Value

	if allowPrivate == allowPublic {
		return allowPrivate
	}

	user, _ := li.Dependencies.Auth.UserOf(r)
	if user == nil {
		return allowPublic
	} else {
		return allowPrivate
	}
}

var (
	_ component.Cronable = (*ListInstances)(nil)
)

func (li *ListInstances) TaskName() string {
	return "instance list and status"
}

func (li *ListInstances) Cron(ctx context.Context) (err error) {
	{
		names, e := li.getNames(ctx)
		if err == nil {
			li.names.Set(names)
		} else {
			err = e
		}
	}

	{
		infos, e := li.getInfos(ctx)
		if err == nil {
			li.infos.Set(infos)
		} else {
			err = e
		}
	}

	return
}

// getNames returns the names of the given instances
func (li *ListInstances) getNames(ctx context.Context) (map[string]struct{}, error) {
	wissKIs, err := li.Dependencies.Instances.All(ctx)
	if err != nil {
		return nil, err
	}

	names := make(map[string]struct{}, len(wissKIs))
	for _, w := range wissKIs {
		names[w.Slug] = struct{}{}
	}
	return names, nil
}

// getInfos returns the names of the given instances
func (li *ListInstances) getInfos(ctx context.Context) ([]status.WissKI, error) {
	// find all the WissKIs
	wissKIs, err := li.Dependencies.Instances.All(ctx)
	if err != nil {
		return nil, err
	}

	infos := make([]status.WissKI, len(wissKIs))

	// determine their infos
	var eg errgroup.Group
	for i, instance := range wissKIs {
		i := i
		wissKI := instance
		eg.Go(func() (err error) {
			infos[i], err = wissKI.Info().Information(ctx, false)
			return
		})
	}
	eg.Wait()

	// filter them by those that are running and do not have prefixes excluded
	infosF := infos[:0]
	for _, info := range infos {
		if info.NoPrefixes || !info.Running {
			continue
		}
		infosF = append(infosF, info)
	}

	// and return them
	return infos[:len(infosF):len(infosF)], nil
}
