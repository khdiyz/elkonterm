-- +goose Up
-- +goose StatementBegin
CREATE TABLE "categories" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "parent_id" uuid,
    "name" jsonb NOT NULL,
    "type" varchar NOT NULL, 
    "photo" varchar(64),
    "is_top" boolean NOT NULL DEFAULT FALSE,
    "application_areas_img" varchar(64),
    "status" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP DEFAULT now(),
    "updated_at" TIMESTAMP NULL,
    "deleted_at" TIMESTAMP NULL,
    FOREIGN KEY ("parent_id") REFERENCES "categories"(id)
); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "categories";
-- +goose StatementEnd
