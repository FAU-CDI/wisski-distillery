package auth

import (
	"sync"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/gorilla/sessions"
)

type Auth struct {
	component.Base
	Dependencies struct {
		SQL *sql.SQL
	}

	storeOnce sync.Once
	store     sessions.Store
}

var (
	_ component.Routeable = (*Auth)(nil)
)
