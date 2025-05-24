//spellchecker:words next
package next

//spellchecker:words context errors http github wisski distillery internal component auth policy scopes instances server handling ingredient users pkglib httpx
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/users"
	"github.com/tkw1536/pkglib/httpx"
)

type Next struct {
	component.Base
	dependencies struct {
		Auth      *auth.Auth
		Policy    *policy.Policy
		Instances *instances.Instances
		Handleing *handling.Handling
	}
}

var (
	_ component.Routeable = (*Next)(nil)
)

func (next *Next) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/next/",
		Decorator: next.dependencies.Auth.Require(true, scopes.ScopeUserValid, nil),
	}
}

// Next returns a url that will forward authorized users to the given slug and path.
func (next *Next) Next(context context.Context, slug, path string) (string, error) {
	wisski, err := next.dependencies.Instances.WissKI(context, slug)
	if err != nil {
		return "", fmt.Errorf("failed to get WissKI: %w", err)
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
	slug, ok := component.GetStill(next).Config.HTTP.SlugFromHost(url.Host)
	if slug == "" || !ok {
		return nil, "", httpx.ErrBadRequest
	}

	// fetch the instance from the database
	wisski, err = next.dependencies.Instances.WissKI(r.Context(), slug)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get WissKI: %w", err)
	}

	// return the wisski and the relative path
	return wisski, url.Path, nil
}

func (next *Next) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return next.dependencies.Handleing.Redirect(func(r *http.Request) (string, int, error) {
		// get the instance and the path
		instance, path, err := next.getInstance(r)
		if err != nil {
			return "", 0, httpx.ErrForbidden
		}

		// get the user
		user, _, err := next.dependencies.Auth.SessionOf(r)
		if err != nil {
			return "", 0, fmt.Errorf("failed to get session: %w", err)
		}

		// check if they have a grant
		grant, err := next.dependencies.Policy.Has(r.Context(), user.User.User, instance.Slug)
		if errors.Is(err, policy.ErrNoAccess) {
			return "", 0, httpx.ErrForbidden
		}
		if err != nil {
			return "", 0, fmt.Errorf("failed to check access: %w", err)
		}

		// perform the login
		dest, err := instance.Users().LoginWithOpt(r.Context(), nil, grant.DrupalUsername, users.LoginOptions{
			Destination:     path,
			CreateIfMissing: true,
			GrantAdminRole:  grant.DrupalAdminRole,
		})
		if err != nil {
			return "", 0, fmt.Errorf("failed to login user: %w", err)
		}

		// and redirect
		return dest.String(), http.StatusSeeOther, nil
	}), nil
}
