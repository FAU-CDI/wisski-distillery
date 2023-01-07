package admin

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/gorilla/mux"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

//go:embed "html/grants.html"
var grantsStr string
var grantsTemplate = static.AssetsAdmin.MustParseShared(
	"grants.html",
	grantsStr,
)

type grantsContext struct {
	custom.BaseContext

	Error string

	instance *wisski.WissKI
	Instance models.Instance // current instance

	Grants    []models.Grant // grants that exist for the user
	Usernames []string       // unuused distillery usernames
	Drupals   []string       // unusued drupal usernames
}

func (gc *grantsContext) use(r *http.Request, slug string, admin *Admin) (err error) {
	admin.Dependencies.Custom.Update(gc, r)

	// find the instance itself
	gc.instance, err = admin.Dependencies.Instances.WissKI(r.Context(), slug)
	if err == instances.ErrWissKINotFound {
		return httpx.ErrNotFound
	}
	if err != nil {
		return err
	}
	gc.Instance = gc.instance.Instance

	return nil
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

func (admin *Admin) getGrants(r *http.Request) (gc grantsContext, err error) {
	if err := gc.use(r, mux.Vars(r)["slug"], admin); err != nil {
		return gc, err
	}

	if err := gc.useGrants(r, admin); err != nil {
		return gc, err
	}

	return gc, nil
}

func (admin *Admin) postGrants(r *http.Request) (gc grantsContext, err error) {
	// parse the form
	if err := r.ParseForm(); err != nil {
		return gc, err
	}

	// read out the form values
	var (
		slug           = r.PostFormValue("slug")
		delete         = r.PostFormValue("action") == "delete"
		distilleryUser = r.PostFormValue("distillery-user")
		drupalUser     = r.PostFormValue("drupal-user")
		adminRole      = r.PostFormValue("admin") == httpx.CheckboxChecked
	)

	// set the common fields
	if err := gc.use(r, slug, admin); err != nil {
		return gc, err
	}

	if delete {
		// delete the user grant
		err := admin.Dependencies.Policy.Remove(r.Context(), distilleryUser, slug)
		if err != nil {
			return gc, err
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
		return gc, err
	}
	return gc, nil
}
