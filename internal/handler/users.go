package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/utilyre/ssa/internal/router"
	"github.com/utilyre/xmate"
)

type User struct {
	ID       int64  `json:"id,omitempty" validate:"isdefault"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=1024"`
}

type usersHandler struct {
	validate *validator.Validate
}

func HandleUsers(r router.Router, v *validator.Validate) {
	s := r.Subrouter("/users")
	h := usersHandler{
		validate: v,
	}

	s.HandleFunc("/signup", h.signup).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")
}

func (h usersHandler) signup(w http.ResponseWriter, r *http.Request) error {
	user := new(User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return xmate.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(user); err != nil {
		return xmate.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return xmate.WriteJSON(w, http.StatusOK, user)
}
