package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func New(lc fx.Lifecycle, l *slog.Logger) *sqlx.DB {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		l.Error("failed to connect to postgres database", "error", err)
		os.Exit(1)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return db.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db
}
