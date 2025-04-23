//spellchecker:words admin
package admin

//spellchecker:words context embed html template http github wisski distillery internal component instances server assets templating models status pkglib httpx form field julienschmidt httprouter golang maps slices
import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/httpx/form/field"

	"maps"
	"slices"

	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance_users.html"
var instanceUsersHTML []byte
var instanceUsersTemplate = templating.Parse[instanceUsersContext](
	"instance_users.html", instanceUsersHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type instanceUsersContext struct {
	templating.RuntimeFlags

	Error string

	instance *wisski.WissKI
	Instance models.Instance // current instance

	Users []status.DrupalUser // drupal users

	Usernames []string       // unuused distillery usernames
	Grants    []models.Grant // grants that exist for the user
}

func (admin *Admin) instanceUsers(context.Context) http.Handler {
	tpl := instanceUsersTemplate.Prepare(
		admin.dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuGrants,
		),
	)

	return tpl.HTMLHandlerWithFlags(admin.dependencies.Handling, func(r *http.Request) (instanceUsersContext, []templating.FlagFunc, error) {
		if r.Method == http.MethodGet {
			return admin.getGrantsUsers(r)
		} else {
			return admin.postInstanceUsers(r)
		}
	})
}

func (admin *Admin) getGrantsUsers(r *http.Request) (gc instanceUsersContext, funcs []templating.FlagFunc, err error) {
	slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

	funcs, err = gc.use(r, slug, admin)
	if err != nil {
		return gc, nil, err
	}

	if err := gc.useUsers(r, admin); err != nil {
		return gc, nil, err
	}

	return gc, funcs, nil
}

func (admin *Admin) postInstanceUsers(r *http.Request) (gc instanceUsersContext, funcs []templating.FlagFunc, err error) {
	// parse the form
	if err := r.ParseForm(); err != nil {
		return gc, nil, err
	}

	// read out the form values
	var (
		slug           = r.PostFormValue("slug")
		actionIsDelete = r.PostFormValue("action") == "delete"
		distilleryUser = r.PostFormValue("distillery-user")
		drupalUser     = r.PostFormValue("drupal-user")
		adminRole      = r.PostFormValue("admin") == field.CheckboxChecked
	)

	// set the common fields
	funcs, err = gc.use(r, slug, admin)
	if err != nil {
		return gc, nil, err
	}

	if actionIsDelete {
		// delete the user grant
		err := admin.dependencies.Policy.Remove(r.Context(), distilleryUser, slug)
		if err != nil {
			return gc, nil, err
		}
	} else {
		// update the grant
		err := admin.dependencies.Policy.Set(r.Context(), models.Grant{
			User: distilleryUser,
			Slug: slug,

			DrupalUsername:  drupalUser,
			DrupalAdminRole: adminRole,
		})
		if err != nil {
			gc.Error = fmt.Sprintf("Unable to update grant for user %s: %s", distilleryUser, err.Error())
		}
	}

	// fetch the grants for the instance
	if err := gc.useUsers(r, admin); err != nil {
		return gc, nil, err
	}
	return gc, funcs, nil
}

func (gc *instanceUsersContext) use(r *http.Request, slug string, admin *Admin) (funcs []templating.FlagFunc, err error) {
	// find the instance itself
	gc.instance, err = admin.dependencies.Instances.WissKI(r.Context(), slug)
	if errors.Is(err, instances.ErrWissKINotFound) {
		return nil, httpx.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	gc.Instance = gc.instance.Instance

	// replace the functions
	escapedSlug := url.PathEscape(slug)
	funcs = []templating.FlagFunc{
		templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + escapedSlug)}),                // #nosec G203 -- escaped and safe
		templating.ReplaceCrumb(menuGrants, component.MenuItem{Title: "Users & Grants", Path: template.URL("/admin/instance/" + escapedSlug + "/users")}), // #nosec G203 -- escaped and safe
		templating.Title(gc.Instance.Slug + " - Users & Grants"),
		admin.instanceTabs(escapedSlug, "users"),
	}
	return funcs, nil
}

func (gc *instanceUsersContext) useUsers(r *http.Request, admin *Admin) (err error) {
	gc.Grants, err = admin.dependencies.Policy.Instance(r.Context(), gc.Instance.Slug)
	if err != nil {
		return err
	}

	users, err := admin.dependencies.Auth.Users(r.Context())
	if err != nil {
		return err
	}

	// create a namemap of users, but not those already taken
	userNameMap := make(map[string]struct{}, len(users))
	for _, user := range users {
		userNameMap[user.User.User] = struct{}{}
	}
	for _, grant := range gc.Grants {
		delete(userNameMap, grant.User)
	}

	// setup the usernames
	gc.Usernames = slices.AppendSeq(make([]string, 0, len(userNameMap)), maps.Keys(userNameMap))
	slices.Sort(gc.Usernames)

	// get the drupal user data
	gc.Users, err = gc.instance.Users().All(r.Context(), nil)
	if err != nil {
		return err
	}
	slices.SortFunc(gc.Users, func(a, b status.DrupalUser) int {
		return int(a.UID) - int(b.UID)
	})

	return nil
}
