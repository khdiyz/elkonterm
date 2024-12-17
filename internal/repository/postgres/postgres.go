package postgres

import (
	"elkonterm/config"
	"elkonterm/pkg/logger"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// SQLXConfig holds the configuration for the database instance
type SQLXConfig struct {
	DriverName     string
	DataSourceName string
	MaxOpenConns   int
	MaxIdleConns   int
	MaxLifetime    time.Duration
}

// InitializeSQLXDatabase returns a new DBInstance
func (cfg *SQLXConfig) InitializeSQLXDatabase(logger *logger.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	logger.Infof("Setting maximum number of open connections to %d", cfg.MaxOpenConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	logger.Infof("Setting maximum number of idle connections to %d", cfg.MaxIdleConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	logger.Infof("Setting maximum lifetime for a connection to %s", cfg.MaxLifetime)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}

// SetupPostgresConnection sets up a connection to the PostgreSQL database using sqlx
func NewPostgresDB(cfg *config.Config, logger *logger.Logger) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable timezone=Asia/Tashkent",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
	)

	// Setup sqlx config for PostgreSQL
	sqlxConfig := SQLXConfig{
		DriverName:     "postgres",
		DataSourceName: dsn,
		MaxOpenConns:   100,
		MaxIdleConns:   10,
		MaxLifetime:    15 * time.Minute,
	}

	// Initialize PostgreSQL connection with sqlx
	db, err := sqlxConfig.InitializeSQLXDatabase(logger)
	if err != nil {
		return nil, err
	}

	return db, nil
}
