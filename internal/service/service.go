package service

import (
	"elkonterm/config"
	"elkonterm/internal/models"
	"elkonterm/internal/repository"
	"elkonterm/pkg/logger"
	"time"
)

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository, cfg *config.Config, loggers *logger.Logger) *Service {
	return &Service{
		Authorization: newAuthService(repos, loggers, cfg),
	}
}

type Authorization interface {
	CreateToken(user models.User, tokenType string, expiresAt time.Time) (*models.Token, error)
	GenerateTokens(user models.User) (*models.Token, *models.Token, error)
	ParseToken(token string) (*jwtCustomClaim, error)
	LoginAdmin(input models.LoginRequest) (*models.Token, *models.Token, error)
}
