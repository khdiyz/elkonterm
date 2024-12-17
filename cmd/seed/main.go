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
	roleSuperAdmin := "SUPER ADMIN"
	roleManager := "MANAGER"
	roleClient := "CLIENT"

	createRoleQuery := `
	INSERT INTO roles (
		id, 
		name
	) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING;`

	_, err = db.Exec(createRoleQuery, config.SuperAdminRoleID, roleSuperAdmin)
	if err != nil {
		logger.Fatalf("Failed to create super admin role: %v\n", err)
	}

	_, err = db.Exec(createRoleQuery, config.ManagerRoleID, roleManager)
	if err != nil {
		logger.Fatalf("Failed to create manager role: %v\n", err)
	}

	_, err = db.Exec(createRoleQuery, config.ClientRoleID, roleClient)
	if err != nil {
		logger.Fatalf("Failed to create client role: %v\n", err)
	}

	// Create New User Super Admin

	userID := config.SuperAdminUserId
	userFullName := "Super Admin"
	userPhone := "+998901234567"
	userEmail := "superadmin@mail.ru"

	userPassword, err := helper.GenerateHash(cfg, "superadmin")
	if err != nil {
		logger.Fatal(err)
	}

	createSuperAdminQuery := `
	INSERT INTO users (
		id,	
		full_name,
		phone_number,
		role_id,
		email,
		password
	) VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (email) DO NOTHING;`

	_, err = db.Exec(createSuperAdminQuery, userID, userFullName, userPhone, config.SuperAdminRoleID, userEmail, userPassword)
	if err != nil {
		logger.Fatalf("Failed to create super admin role: %v\n", err)
	}

	logger.Info("Seed completed successfully!")
}
