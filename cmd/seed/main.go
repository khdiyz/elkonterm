package main

import (
	"elkonterm/config"
	"elkonterm/internal/repository/postgres"
	"elkonterm/pkg/helper"
	"elkonterm/pkg/logger"
)

func main() {
	// Load configuration and logger
	cfg := config.GetConfig()
	logger := logger.GetLogger()

	// Setup PostgreSQL connection
	db, err := postgres.NewPostgresDB(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer db.Close()

	// Create New Roles
	roleAdmin := "ADMIN"
	roleManager := "MANAGER"
	roleClient := "CLIENT"

	createRoleQuery := `
	INSERT INTO roles (
		id, 
		name
	) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING;`

	_, err = db.Exec(createRoleQuery, config.AdminRoleID, roleAdmin)
	if err != nil {
		logger.Fatalf("Failed to create admin role: %v\n", err)
	}

	_, err = db.Exec(createRoleQuery, config.ManagerRoleID, roleManager)
	if err != nil {
		logger.Fatalf("Failed to create manager role: %v\n", err)
	}

	_, err = db.Exec(createRoleQuery, config.ClientRoleID, roleClient)
	if err != nil {
		logger.Fatalf("Failed to create client role: %v\n", err)
	}

	// Create New User Admin

	userID := config.AdminUserId
	userFullName := "Admin"
	userPhone := "+998901234567"
	userEmail := "admin@mail.ru"

	userPassword, err := helper.GenerateHash(cfg, "admin")
	if err != nil {
		logger.Fatal(err)
	}

	createAdminQuery := `
	INSERT INTO users (
		id,	
		full_name,
		phone_number,
		role_id,
		email,
		password
	) VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (email) DO NOTHING;`

	_, err = db.Exec(createAdminQuery, userID, userFullName, userPhone, config.AdminRoleID, userEmail, userPassword)
	if err != nil {
		logger.Fatalf("Failed to create admin role: %v\n", err)
	}

	logger.Info("Seed completed successfully!")
}
