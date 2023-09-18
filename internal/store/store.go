package store

import (
	"os"

	"github.com/gorilla/sessions"
)

func New() sessions.Store {
	store := sessions.NewCookieStore([]byte(os.Getenv("AUTH_KEY")))
	store.Options.HttpOnly = true

	return store
}
