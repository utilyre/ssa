package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Session struct {
	ID        int32     `db:"id"`
	UUID      uuid.UUID `db:"uuid"`
	CreatedAt time.Time `db:"created_at"`

	UserID int32 `db:"user_id"`
	User
}

type SessionStorage struct {
	db *sqlx.DB
}

func NewSessionStorage(db *sqlx.DB) SessionStorage {
	return SessionStorage{db: db}
}

func (s SessionStorage) Create(ctx context.Context, session *Session) error {
	query := `
	INSERT
	INTO "sessions"
	("user_id")
	VALUES
		($1)
	RETURNING "id", "uuid", "created_at";
	`

	return s.db.GetContext(ctx, session, query, session.UserID)
}
