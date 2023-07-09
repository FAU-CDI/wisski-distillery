package admin

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/httpx"
	"github.com/tkw1536/pkglib/httpx/field"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

//go:embed "html/grants.html"
var grantsHTML []byte
var grantsTemplate = templating.Parse[grantsContext](
	"grants.html", grantsHTML, nil,

	templating.Assets(assets.AssetsAdmin),
)

type grantsContext struct {
	templating.RuntimeFlags

	Error string

	instance *wisski.WissKI
	Instance models.Instance // current instance

	Grants    []models.Grant // grants that exist for the user
	Usernames []string       // unuused distillery usernames
	Drupals   []string       // unusued drupal usernames
}

func (admin *Admin) grants(ctx context.Context) http.Handler {
	tpl := grantsTemplate.Prepare(
		admin.Dependencies.Templating,
		templating.Crumbs(
			menuAdmin,
			menuInstances,
			menuInstance,
			menuGrants,
		),
	)

	return tpl.HTMLHandlerWithFlags(func(r *http.Request) (grantsContext, []templating.FlagFunc, error) {
		if r.Method == http.MethodGet {
			return admin.getGrants(r)
		} else {
			return admin.postGrants(r)
		}
	})
}

func (admin *Admin) getGrants(r *http.Request) (gc grantsContext, funcs []templating.FlagFunc, err error) {
	slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")

	funcs, err = gc.use(r, slug, admin)
	if err != nil {
		return gc, nil, err
	}

	if err := gc.useGrants(r, admin); err != nil {
		return gc, nil, err
	}

	return gc, funcs, nil
}

func (admin *Admin) postGrants(r *http.Request) (gc grantsContext, funcs []templating.FlagFunc, err error) {
	// parse the form
	if err := r.ParseForm(); err != nil {
		return gc, nil, err
	}

	// read out the form values
	var (
		slug           = r.PostFormValue("slug")
		delete         = r.PostFormValue("action") == "delete"
		distilleryUser = r.PostFormValue("distillery-user")
		drupalUser     = r.PostFormValue("drupal-user")
		adminRole      = r.PostFormValue("admin") == field.CheckboxChecked
	)

	// set the common fields
	funcs, err = gc.use(r, slug, admin)
	if err != nil {
		return gc, nil, err
	}

	if delete {
		// delete the user grant
		err := admin.Dependencies.Policy.Remove(r.Context(), distilleryUser, slug)
		if err != nil {
			return gc, nil, err
		}
	} else {
		// update the grant
		err := admin.Dependencies.Policy.Set(r.Context(), models.Grant{
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
	if err := gc.useGrants(r, admin); err != nil {
		return gc, nil, err
	}
	return gc, funcs, nil
}

func (gc *grantsContext) use(r *http.Request, slug string, admin *Admin) (funcs []templating.FlagFunc, err error) {
	// find the instance itself
	gc.instance, err = admin.Dependencies.Instances.WissKI(r.Context(), slug)
	if err == instances.ErrWissKINotFound {
		return nil, httpx.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	gc.Instance = gc.instance.Instance

	// replace the functions
	funcs = []templating.FlagFunc{
		templating.ReplaceCrumb(menuInstance, component.MenuItem{Title: "Instance", Path: template.URL("/admin/instance/" + slug)}),
		templating.ReplaceCrumb(menuGrants, component.MenuItem{Title: "Grants", Path: template.URL("/admin/grants/" + slug)}),
		templating.Title(gc.Instance.Slug + " - Grants"),
	}
	return funcs, nil
}

func (gc *grantsContext) useGrants(r *http.Request, admin *Admin) (err error) {
	gc.Grants, err = admin.Dependencies.Policy.Instance(r.Context(), gc.Instance.Slug)
	if err != nil {
		return err
	}

	users, err := admin.Dependencies.Auth.Users(r.Context())
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
	gc.Usernames = maps.Keys(userNameMap)
	slices.Sort(gc.Usernames)

	// get the drupal usernames
	drupals, err := gc.instance.Users().All(r.Context(), nil)
	if err != nil {
		return err
	}

	// and convert them to strings only
	gc.Drupals = make([]string, len(drupals))
	for i, drupal := range drupals {
		gc.Drupals[i] = string(drupal.Name)
	}
	slices.Sort(gc.Drupals)

	return nil
}
