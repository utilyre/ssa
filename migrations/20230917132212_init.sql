-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
    "id" serial NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    "email" varchar(255) NOT NULL UNIQUE,
    "password" bytea NOT NULL,

    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "users";
-- +goose StatementEnd
