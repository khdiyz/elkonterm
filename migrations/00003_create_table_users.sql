-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "full_name" VARCHAR(255) NOT NULL,
    "phone_number" VARCHAR(50),
    "role_id" uuid NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "company" jsonb,
    "status" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP NULL,
    "deleted_at" TIMESTAMP NULL,
    FOREIGN KEY ("role_id") REFERENCES "roles"(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "users";
-- +goose StatementEnd
