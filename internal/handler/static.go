package handler

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/utilyre/ssa/internal/router"
	"github.com/utilyre/ssa/internal/templates"
	"github.com/utilyre/xmate"
)

type staticHandler struct {
	store sessions.Store
	tmpls templates.Templates
}

func HandleStatic(
	r router.Router,
	store sessions.Store,
	t templates.Templates,
) {
	h := staticHandler{
		store: store,
		tmpls: t,
	}

	r.HandleFunc("/", h.home).Methods(http.MethodGet)
	r.HandleFunc("/signup", h.signup).Methods(http.MethodGet)
	r.HandleFunc("/login", h.login).Methods(http.MethodGet)
}

func (h staticHandler) home(w http.ResponseWriter, r *http.Request) error {
	session, err := h.store.Get(r, "ssa-login")
	if err != nil {
		return err
	}

	return xmate.WriteHTML(w, h.tmpls.Pages, http.StatusOK, xmate.Map{
		"Name": "home",
		"Payload": xmate.Map{
			"Email": session.Values["email"],
		},
	})
}

func (h staticHandler) signup(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, h.tmpls.Pages, http.StatusOK, xmate.Map{
		"Name":    "signup",
		"Payload": nil,
	})
}

func (h staticHandler) login(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, h.tmpls.Pages, http.StatusOK, xmate.Map{
		"Name":    "login",
		"Payload": nil,
	})
}
