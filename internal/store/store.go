package store

import (
	"github.com/gorilla/sessions"
	"github.com/utilyre/ssa/internal/env"
)

func New(e env.Env) sessions.Store {
	store := sessions.NewCookieStore(e.AuthKey)
	store.Options.HttpOnly = true

	return store
}
