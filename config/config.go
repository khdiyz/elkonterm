package config

import (
	"elkonterm/pkg/logger"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	HTTPHost string
	HTTPPort int

	Environment string
	Debug       bool

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	JWTSecret                string
	JWTAccessExpirationHours int
	JWTRefreshExpirationDays int

	HashKey string

	PaymeMerchantID  string
	PaymeKey         string
	PaymeLogin       string
	PaymeBaseURL     string
	PaymeCallbackURL string
	PaymeParam       string
}

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			HTTPHost:    cast.ToString(getOrReturnDefault("HTTP_HOST", "localhost")),
			HTTPPort:    cast.ToInt(getOrReturnDefault("HTTP_PORT", 7070)),
			Environment: cast.ToString(getOrReturnDefault("ENVIRONMENT", "development")),
			Debug:       cast.ToBool(getOrReturnDefault("DEBUG", true)),

			PostgresHost:     cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost")),
			PostgresPort:     cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432)),
			PostgresDatabase: cast.ToString(getOrReturnDefault("POSTGRES_DB", "")),
			PostgresUser:     cast.ToString(getOrReturnDefault("POSTGRES_USER", "")),
			PostgresPassword: cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "")),

			RedisHost:     cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost")),
			RedisPort:     cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379)),
			RedisPassword: cast.ToString(getOrReturnDefault("REDIS_PASSWORD", "")),
			RedisDB:       cast.ToInt(getOrReturnDefault("REDIS_DB", 0)),

			JWTSecret:                cast.ToString(getOrReturnDefault("JWT_SECRET", "elkonterm-forever-2025")),
			JWTAccessExpirationHours: cast.ToInt(getOrReturnDefault("JWT_ACCESS_EXPIRATION_HOURS", 12)),
			JWTRefreshExpirationDays: cast.ToInt(getOrReturnDefault("JWT_REFRESH_EXPIRATION_DAYS", 3)),

			HashKey: cast.ToString(getOrReturnDefault("HASH_KEY", "skd32r8wdahkkN2HSdqw")),

			PaymeMerchantID:  cast.ToString(getOrReturnDefault("PAYME_MERCHANT_ID", "")),
			PaymeKey:         cast.ToString(getOrReturnDefault("PAYME_KEY", "")),
			PaymeLogin:       cast.ToString(getOrReturnDefault("PAYME_LOGIN", "")),
			PaymeBaseURL:     cast.ToString(getOrReturnDefault("PAYME_BASE_URL", "")),
			PaymeParam:       cast.ToString(getOrReturnDefault("PAYME_PARAM", "")),
			PaymeCallbackURL: cast.ToString(getOrReturnDefault("PAYME_CALLBACK_URL", "")),
		}
	})

	return instance
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	err := godotenv.Load(".env")
	if err != nil {
		logger.GetLogger().Error(err)
	}
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
