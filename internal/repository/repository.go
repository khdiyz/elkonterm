package repository

import (
	"elkonterm/internal/repository/dto"
	"elkonterm/internal/repository/postgres"
	"elkonterm/pkg/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User
	Category
}

func NewRepository(db *sqlx.DB, logger *logger.Logger) *Repository {
	return &Repository{
		User:     postgres.NewUserRepo(db, logger),
		Category: postgres.NewCategoryRepo(db, logger),
	}
}

type User interface {
	Create(input dto.CreateUser) (uuid.UUID, error)
	GetList(options dto.FilterOptions) ([]dto.User, int, error)
	GetByEmail(email string) (*dto.User, error)
	GetByID(id uuid.UUID) (*dto.User, error)
}

type Category interface {
	Create(input dto.CreateCategory) (uuid.UUID, error)
	GetList(options dto.FilterOptions) ([]dto.Category, int, error)
}
