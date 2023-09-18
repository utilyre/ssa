package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/utilyre/ssa/internal/router"
	"github.com/utilyre/ssa/internal/storage"
	"github.com/utilyre/xmate"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id,omitempty" validate:"isdefault"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=1024"`
}

type usersHandler struct {
	storage  storage.UserStorage
	store    sessions.Store
	validate *validator.Validate
}

func HandleUsers(
	r router.Router,
	s storage.UserStorage,
	store sessions.Store,
	v *validator.Validate,
) {
	sr := r.Subrouter("/users")
	h := usersHandler{
		storage:  s,
		store:    store,
		validate: v,
	}

	sr.HandleFunc("/signup", h.signup).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json")

	sr.HandleFunc("/login", h.login).
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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	dbUser := &storage.User{
		Email:    user.Email,
		Password: hash,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1800*time.Millisecond)
	defer cancel()

	if err := h.storage.Create(ctx, dbUser); err != nil {
		if errors.Is(err, storage.ErrDuplicateKey) {
			return xmate.NewHTTPError(http.StatusConflict, "user already exists")
		}

		return err
	}

	user.ID = dbUser.ID

	w.Header().Set("HX-Redirect", "/login")
	return xmate.WriteJSON(w, http.StatusCreated, user)
}

func (h usersHandler) login(w http.ResponseWriter, r *http.Request) error {
	user := new(User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return xmate.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(user); err != nil {
		return xmate.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dbUser := &storage.User{
		Email: user.Email,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1800*time.Millisecond)
	defer cancel()

	if err := h.storage.ReadByEmail(ctx, dbUser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return xmate.NewHTTPError(http.StatusNotFound)
		}

		return err
	}

	if err := bcrypt.CompareHashAndPassword(dbUser.Password, []byte(user.Password)); err != nil {
		return xmate.NewHTTPError(http.StatusNotFound)
	}

	session, err := h.store.Get(r, "ssa-login")
	if err != nil {
		return err
	}

	session.Values["email"] = user.Email

	w.Header().Set("HX-Redirect", "/")
	return session.Save(r, w)
}
