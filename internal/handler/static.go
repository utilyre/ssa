package handler

import (
	"html/template"
	"net/http"

	"github.com/utilyre/ssa/internal/router"
	"github.com/utilyre/xmate"
)

type staticHandler struct {
	tmpl *template.Template
}

func HandleStatic(r router.Router, t *template.Template) {
	h := staticHandler{
		tmpl: t,
	}

	r.Handle("/signup", r.ErrorHandler.HandleFunc(h.signup)).
		Methods(http.MethodGet)

	r.Handle("/login", r.ErrorHandler.HandleFunc(h.login)).
		Methods(http.MethodGet)
}

func (h staticHandler) signup(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, h.tmpl, http.StatusOK, "signup", nil)
}

func (h staticHandler) login(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, h.tmpl, http.StatusOK, "login", nil)
}
