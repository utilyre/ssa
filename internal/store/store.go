package store

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/utilyre/ssa/internal/storage"
)

type UUIDKey struct{}
type UserIDKey struct{}

type dbStore struct {
	sc      *securecookie.SecureCookie
	storage storage.SessionStorage
}

func New(s storage.SessionStorage, l *slog.Logger) sessions.Store {
	hashKey, err := hex.DecodeString(os.Getenv("HASH_KEY"))
	if err != nil {
		l.Error("failed to decode HASH_KEY environment variable", "error", err)
		os.Exit(1)
	}

	blockKey, err := hex.DecodeString(os.Getenv("BLOCK_KEY"))
	if err != nil {
		l.Error("failed to decode BLOCK_KEY environment variable", "error", err)
		os.Exit(1)
	}

	return dbStore{
		sc:      securecookie.New(hashKey, blockKey),
		storage: s,
	}
}

func (s dbStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	session, err := s.New(r, name)
	if err != nil {
		return nil, err
	}

	cookie, err := r.Cookie(name)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return session, nil
		default:
			return nil, err
		}
	}

	var id uuid.UUID
	if err := s.sc.Decode(name, cookie.Value, &id); err != nil {
		return nil, err
	}

	session.Values[UUIDKey{}] = id
	return session, nil
}

func (s dbStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.NewSession(s, name), nil
}

func (s dbStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	anyUUID, ok := session.Values[UUIDKey{}]
	userID := session.Values[UserIDKey{}].(int32)

	var id uuid.UUID
	if ok {
		id = anyUUID.(uuid.UUID)
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 1800*time.Millisecond)
		defer cancel()

		dbSession := &storage.Session{UserID: userID}
		s.storage.Create(ctx, dbSession)
		id = dbSession.UUID
	}

	encoded, err := s.sc.Encode(session.Name(), id)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  session.Name(),
		Value: encoded,
		Path:  "/",
		// Secure: true,
		HttpOnly: true,
	})
	return nil
}
