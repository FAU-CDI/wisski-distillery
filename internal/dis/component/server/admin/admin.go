package admin

import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/admin/socket"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/julienschmidt/httprouter"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/tkw1536/pkglib/httpx"
)

type Admin struct {
	component.Base
	dependencies struct {
		Handling *handling.Handling
		Fetchers []component.DistilleryFetcher

		Instances *instances.Instances

		Auth *auth.Auth

		Policy *policy.Policy

		Templating *templating.Templating

		Sockets *socket.Sockets
	}
}

var (
	_ component.DistilleryFetcher = (*Admin)(nil)
	_ component.Routeable         = (*Admin)(nil)
	_ component.Menuable          = (*Admin)(nil)
)

func (admin *Admin) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/admin/",
		CSRF:      true,
		Decorator: admin.dependencies.Auth.Require(false, scopes.ScopeUserAdmin, nil),
	}
}

func (admin *Admin) Menu(r *http.Request) []component.MenuItem {
	if admin.dependencies.Auth.CheckScope("", scopes.ScopeUserAdmin, r) != nil {
		return nil
	}
	return []component.MenuItem{
		{
			Title:    "Admin",
			Path:     "/admin/",
			Priority: component.MenuAdmin,
		},
	}
}

var (
	menuAdmin = component.MenuItem{Title: "Admin", Path: "/admin/"}

	menuUsers      = component.MenuItem{Title: "Users", Path: "/admin/users/"}
	menuUserCreate = component.MenuItem{Title: "Create User", Path: "/admin/users/create/"}

	menuProvision = component.MenuItem{Title: "Provision", Path: "/admin/instances/provision/"}

	menuInstances   = component.MenuItem{Title: "Instances", Path: "/admin/instances/"}
	menuInstance    = component.DummyMenuItem()
	menuRebuild     = component.DummyMenuItem()
	menuGrants      = component.DummyMenuItem()
	menuPurge       = component.DummyMenuItem()
	menuSnapshots   = component.DummyMenuItem()
	menuSSH         = component.DummyMenuItem()
	menuStats       = component.DummyMenuItem()
	menuData        = component.DummyMenuItem()
	menuTriplestore = component.DummyMenuItem()
	menuDrupal      = component.DummyMenuItem()
)

func (admin *Admin) HandleRoute(ctx context.Context, route string) (handler http.Handler, err error) {

	router := httprouter.New()

	// add a handler for the index page
	{
		index := admin.index(ctx)
		router.Handler(http.MethodGet, route, index)
	}

	// add a handler for the instances (and provision) page
	{
		instances := admin.instances(ctx)
		router.Handler(http.MethodGet, route+"instances/", instances)

		provision := admin.instanceProvision(ctx)
		router.Handler(http.MethodGet, route+"instances/provision", provision)
	}

	// add a handler for the user page
	{
		users := admin.users(ctx)
		router.Handler(http.MethodGet, route+"users", users)
	}

	// add a user create form
	{
		create := admin.createUser(ctx)
		router.Handler(http.MethodGet, route+"users/create", create)
		router.Handler(http.MethodPost, route+"users/create", create)
	}

	// add all the admin actions
	router.Handler(http.MethodPost, route+"users/delete", admin.usersDeleteHandler(ctx))
	router.Handler(http.MethodPost, route+"users/disable", admin.usersDisableHandler(ctx))
	router.Handler(http.MethodPost, route+"users/disabletotp", admin.usersDisableTOTPHandler(ctx))
	router.Handler(http.MethodPost, route+"users/password", admin.usersPasswordHandler(ctx))
	router.Handler(http.MethodPost, route+"users/toggleadmin", admin.usersToggleAdmin(ctx))
	router.Handler(http.MethodPost, route+"users/impersonate", admin.usersImpersonateHandler(ctx))
	router.Handler(http.MethodPost, route+"users/unsetpassword", admin.usersUnsetPasswordHandler(ctx))

	// add a handler for the instance page
	{
		instance := admin.instance(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug", instance)
	}

	{
		rebuild := admin.instanceRebuild(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug/rebuild", rebuild)
	}

	{
		iUsers := admin.instanceUsers(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug/users", iUsers)
		router.Handler(http.MethodPost, route+"grants/", iUsers) // NOTE(twiesing): This path is intentionally different!
	}

	{
		purge := admin.instancePurge(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug/purge", purge)
	}

	{
		snapshots := admin.instanceSnapshots(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug/snapshots", snapshots)
	}

	{
		ssh := admin.instanceSSH(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug/ssh", ssh)
	}

	{
		stats := admin.instanceStats(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug/stats", stats)
	}

	{
		triplestore := admin.instanceTS(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug/triplestore", triplestore)
	}

	{
		data := admin.instanceData(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug/data", data)
	}

	{
		drupal := admin.instanceDrupal(ctx)
		router.Handler(http.MethodGet, route+"instance/:slug/drupal", drupal)
	}

	// add a router for the login page
	router.Handler(http.MethodPost, route+"login", admin.loginHandler(ctx))

	return router, nil
}

func (admin *Admin) loginHandler(ctx context.Context) http.Handler {
	logger := wdlog.Of(ctx)

	return admin.dependencies.Handling.Redirect(func(r *http.Request) (string, int, error) {
		// parse the form
		if err := r.ParseForm(); err != nil {
			logger.Error(
				"failed to parse admin login",
				"error", err,
			)
			return "", 0, err
		}

		// get the instance
		instance, err := admin.dependencies.Instances.WissKI(r.Context(), r.PostFormValue("slug"))
		if err != nil {
			return "", 0, httpx.ErrNotFound
		}

		target, err := instance.Users().Login(r.Context(), nil, r.PostFormValue("user"))
		if err != nil {
			logger.Error(
				"failed to admin login",
				"error", err,
			)
			return "", 0, err
		}
		return target.String(), http.StatusSeeOther, err
	})
}
