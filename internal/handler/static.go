package handler

import (
	"net/http"

	"github.com/utilyre/ssa/internal/router"
	"github.com/utilyre/ssa/internal/templates"
	"github.com/utilyre/xmate"
)

type staticHandler struct {
	tmpls templates.Templates
}

func HandleStatic(r router.Router, t templates.Templates) {
	h := staticHandler{
		tmpls: t,
	}

	r.HandleFunc("/signup", h.signup).Methods(http.MethodGet)
	r.HandleFunc("/login", h.login).Methods(http.MethodGet)
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
