-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION "uuid-ossp";

CREATE TABLE "sessions" (
    "id" serial NOT NULL,
    "uuid" uuid NOT NULL DEFAULT uuid_generate_v1(),
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    "user_id" int4 NOT NULL REFERENCES users("id"), 
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "sessions";
DROP EXTENSION "uuid-ossp";
-- +goose StatementEnd
