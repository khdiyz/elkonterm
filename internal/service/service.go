package service

import (
	"elkonterm/config"
	"elkonterm/internal/models"
	"elkonterm/internal/repository"
	"elkonterm/internal/repository/dto"
	"elkonterm/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	Authorization
	Category
}

func NewService(repos *repository.Repository, cfg *config.Config, loggers *logger.Logger) *Service {
	return &Service{
		Authorization: newAuthService(repos, loggers, cfg),
		Category:      newCategoryService(repos, loggers),
	}
}

type Authorization interface {
	CreateToken(user models.User, tokenType string, expiresAt time.Time) (*models.Token, error)
	GenerateTokens(user models.User) (*models.Token, *models.Token, error)
	ParseToken(token string) (*jwtCustomClaim, error)
	LoginAdmin(input models.LoginRequest) (*models.Token, *models.Token, error)
}

type Category interface {
	CreateCategory(input models.CreateCategory) (uuid.UUID, error)
	GetList(options dto.FilterOptions) ([]models.Category, int, error)
}
