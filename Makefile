tidy:
	@go mod tidy

run:
	@go run cmd/app/main.go

up:
	@go run cmd/migration/main.go up

down:
	@go run cmd/migration/main.go down

redo:
	@go run cmd/migration/main.go redo

status:
	@go run cmd/migration/main.go status

create:
	@read -p "Enter migration name: " MIGRATION_NAME; \
	echo "Creating migration: $$MIGRATION_NAME"; \
	go run cmd/migration/main.go create $$MIGRATION_NAME

seed:
	@go run cmd/seed/main.go

swag:
	@swag init -g cmd/app/main.go