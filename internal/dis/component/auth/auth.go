package auth

import (
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/gorilla/sessions"
)

type Auth struct {
	component.Base
	Dependencies struct {
		SQL *sql.SQL
	}

	store lazy.Lazy[sessions.Store]
	csrf  lazy.Lazy[func(http.Handler) http.Handler]
}

var (
	_ component.Routeable = (*Auth)(nil)
)
