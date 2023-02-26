package next

import (
	"context"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/users"
	"github.com/tkw1536/pkglib/httpx"
)

type Next struct {
	component.Base
	Dependencies struct {
		Auth      *auth.Auth
		Policy    *policy.Policy
		Instances *instances.Instances
	}
}

var (
	_ component.Routeable = (*Next)(nil)
)

func (next *Next) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/next/",
		Decorator: next.Dependencies.Auth.Require(auth.User),
	}
}

// Next returns a url that will forward authorized users to the given slug and path
func (next *Next) Next(context context.Context, slug, path string) (string, error) {
	wisski, err := next.Dependencies.Instances.WissKI(context, slug)
	if err != nil {
		return "", err
	}

	target := wisski.URL()
	target.Path = path
	return "/next/?next=" + url.PathEscape(target.String()), nil

}

func (next *Next) getInstance(r *http.Request) (wisski *wisski.WissKI, path string, err error) {
	// extract the instance
	url, err := url.Parse(r.URL.Query().Get("next"))
	if err != nil {
		return nil, "", httpx.ErrBadRequest
	}

	// find the slug
	slug, ok := next.Config.HTTP.SlugFromHost(url.Host)
	if slug == "" || !ok {
		return nil, "", httpx.ErrBadRequest
	}

	// fetch the instance from the database
	wisski, err = next.Dependencies.Instances.WissKI(r.Context(), slug)
	if err != nil {
		return nil, "", err
	}

	// return the wisski and the relative path
	return wisski, url.Path, nil
}

func (next *Next) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return httpx.RedirectHandler(func(r *http.Request) (string, int, error) {
		// get the instance and the path
		instance, path, err := next.getInstance(r)
		if err != nil {
			return "", 0, httpx.ErrForbidden
		}

		// get the user
		user, err := next.Dependencies.Auth.UserOf(r)
		if err != nil {
			return "", 0, err
		}

		// check if they have a grant
		grant, err := next.Dependencies.Policy.Has(r.Context(), user.User.User, instance.Slug)
		if err == policy.ErrNoAccess {
			return "", 0, httpx.ErrForbidden
		}
		if err != nil {
			return "", 0, err
		}

		// perform the login
		dest, err := instance.Users().LoginWithOpt(r.Context(), nil, grant.DrupalUsername, users.LoginOptions{
			Destination:     path,
			CreateIfMissing: true,
			GrantAdminRole:  grant.DrupalAdminRole,
		})
		if err != nil {
			return "", 0, err
		}

		// and redirect
		return dest.String(), http.StatusSeeOther, nil
	}), nil
}
