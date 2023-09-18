package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var ErrDuplicateKey = errors.New("duplicate key value violates unique constraint")

type User struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`

	Email    string `db:"email"`
	Password []byte `db:"password"`
}

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) UserStorage {
	return UserStorage{db: db}
}

func (s UserStorage) Create(ctx context.Context, user *User) error {
	query := `
	INSERT
	INTO "users"
	("email", "password")
	VALUES
		($1, $2)
	RETURNING "id", "created_at";
	`

	if err := s.db.GetContext(ctx, user, query, user.Email, user.Password); err != nil {
		pqErr := new(pq.Error)
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrDuplicateKey
		}

		return err
	}

	return nil
}
