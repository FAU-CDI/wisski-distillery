//spellchecker:words panel
package panel

//spellchecker:words context http github wisski distillery internal component auth next policy scopes tokens instances server handling templating sshkeys models julienschmidt httprouter pkglib httpx form
import (
	"context"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/next"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/policy"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/tokens"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/ssh2/sshkeys"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/julienschmidt/httprouter"
	"go.tkw01536.de/pkglib/httpx/form"
)

type UserPanel struct {
	component.Base
	dependencies struct {
		Auth *auth.Auth

		Handling   *handling.Handling
		Templating *templating.Templating

		Policy    *policy.Policy
		Tokens    *tokens.Tokens
		Instances *instances.Instances
		Next      *next.Next
		Keys      *sshkeys.SSHKeys
		SSH2      *ssh2.SSH2
	}
}

var (
	_ component.Routeable = (*UserPanel)(nil)
	_ component.Menuable  = (*UserPanel)(nil)
)

func (panel *UserPanel) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/user/",
		CSRF:      true,
		Decorator: panel.dependencies.Auth.Require(false, scopes.ScopeUserValid, nil),
	}
}

func (panel *UserPanel) Menu(r *http.Request) []component.MenuItem {
	title := "Login"

	user, err := panel.dependencies.Auth.UserOfSession(r)
	if user != nil && err == nil {
		title = user.User.User
	}
	return []component.MenuItem{
		{Title: title, Priority: component.MenuUser, Path: "/user/"},
	}
}

var (
	menuUser           = component.MenuItem{Title: "User", Path: "/user/"}
	menuChangePassword = component.MenuItem{Title: "Change Password", Path: "/user/password/"}
	menuSSH            = component.MenuItem{Title: "SSH Keys", Path: "/user/ssh/"}
	menuSSHAdd         = component.MenuItem{Title: "Add New Key", Path: "/user/ssh/add/"}

	menuTokens    = component.MenuItem{Title: "Tokens", Path: "/user/tokens/"}
	menuTokensAdd = component.MenuItem{Title: "Add New Token", Path: "/user/tokens/add/"}

	menuTOTPAction  = component.DummyMenuItem()
	menuTOTPDisable = component.MenuItem{Title: "Disable Passcode (TOTP)", Path: "/user/totp/disable/"}
	menuTOTPEnable  = component.MenuItem{Title: "Enable Passcode (TOTP)", Path: "/user/totp/enable/"}
)

func (panel *UserPanel) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	router := httprouter.New()

	{
		user := panel.routeUser(ctx)
		router.Handler(http.MethodGet, route, user)
	}

	{
		password := panel.routePassword(ctx)
		router.Handler(http.MethodGet, route+"password", password)
		router.Handler(http.MethodPost, route+"password", password)
	}

	{
		totpEnable := panel.routeTOTPEnable(ctx)
		router.Handler(http.MethodGet, route+"totp/enable", totpEnable)
		router.Handler(http.MethodPost, route+"totp/enable", totpEnable)
	}

	{
		totpEnroll := panel.routeTOTPEnroll(ctx)
		router.Handler(http.MethodGet, route+"totp/enroll", totpEnroll)
		router.Handler(http.MethodPost, route+"totp/enroll", totpEnroll)
	}

	{
		totpdisable := panel.routeTOTPDisable(ctx)
		router.Handler(http.MethodGet, route+"totp/disable", totpdisable)
		router.Handler(http.MethodPost, route+"totp/disable", totpdisable)
	}

	{
		ssh := panel.sshRoute(ctx)
		router.Handler(http.MethodGet, route+"ssh", ssh)
	}

	{
		sshAdd := panel.sshAddRoute(ctx)
		router.Handler(http.MethodGet, route+"ssh/add", sshAdd)
		router.Handler(http.MethodPost, route+"ssh/add", sshAdd)
	}

	{
		sshDelete := panel.sshDeleteRoute(ctx)
		router.Handler(http.MethodPost, route+"ssh/delete", sshDelete)
	}

	{
		tokens := panel.tokensRoute(ctx)
		router.Handler(http.MethodGet, route+"tokens", tokens)
	}

	{
		tokensAdd := panel.tokensAddRoute(ctx)
		router.Handler(http.MethodGet, route+"tokens/add", tokensAdd)
		router.Handler(http.MethodPost, route+"tokens/add", tokensAdd)
	}

	{
		tokensDelete := panel.tokensDeleteRoute(ctx)
		router.Handler(http.MethodPost, route+"tokens/delete", tokensDelete)
	}

	// ensure that the user is logged in!
	return panel.dependencies.Auth.Protect(router, false, scopes.ScopeUserValid, nil), nil
}

//nolint:errname
type userFormContext struct {
	templating.RuntimeFlags
	form.FormContext

	User *models.User
}

func (panel *UserPanel) UserFormContext(tpl *templating.Template[userFormContext], last component.MenuItem, funcs ...templating.FlagFunc) func(ctx form.FormContext, r *http.Request) any {
	funcs = append(funcs, func(flags templating.Flags, r *http.Request) templating.Flags {
		// append the last menu item, and prepend the menuUser one!
		flags.Crumbs = append(flags.Crumbs, last, last)
		copy(flags.Crumbs[1:], flags.Crumbs)
		flags.Crumbs[0] = menuUser
		return flags
	})

	return func(ctx form.FormContext, r *http.Request) any {
		uctx := userFormContext{FormContext: ctx}
		if user, err := panel.dependencies.Auth.UserOfSession(r); err == nil {
			uctx.User = &user.User
		}
		return tpl.Context(r, uctx, funcs...)
	}
}
