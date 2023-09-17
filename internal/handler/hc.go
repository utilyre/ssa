package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/utilyre/ssa/internal/router"
	"github.com/utilyre/xmate"
)

func HandleHC(r router.Router, db *sqlx.DB) {
	r.Handle("/hc", r.ErrorHandler.HandleFunc(handleHC))
}

func handleHC(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteText(w, http.StatusOK, "Everything is good to Go!")
}
