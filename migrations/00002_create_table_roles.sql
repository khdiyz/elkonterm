-- +goose Up
-- +goose StatementBegin
CREATE TABLE "roles" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(64) NOT NULL UNIQUE,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP NULL,
    "deleted_at" TIMESTAMP NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "roles";
-- +goose StatementEnd
