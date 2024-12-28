package service

import (
	"elkonterm/config"
	"elkonterm/internal/models"
	"elkonterm/internal/repository"
	"elkonterm/internal/repository/dto"
	"elkonterm/pkg/logger"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type categoryService struct {
	repo   *repository.Repository
	logger *logger.Logger
}

func newCategoryService(repo *repository.Repository, logger *logger.Logger) *categoryService {
	return &categoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *categoryService) CreateCategory(input models.CreateCategory) (uuid.UUID, error) {
	var parentID *uuid.UUID

	if input.ParentID != "" {
		parsedID, err := uuid.Parse(input.ParentID)
		if err != nil {
			return uuid.Nil, serviceError(err, codes.InvalidArgument)
		}

		parentID = &parsedID
	}

	if input.Type != config.SystemTypeElcon && input.Type != config.SystemTypeFinfire {
		return uuid.Nil, serviceError(errors.New("type must be only: 'elcon' or 'finfire'"), codes.InvalidArgument)
	}

	id, err := s.repo.Category.Create(dto.CreateCategory{
		ParentID:            parentID,
		Name:                input.Name,
		Type:                input.Type,
		Photo:               input.Photo,
		IsTop:               input.IsTop,
		ApplicationAreasImg: input.ApplicationAreasImg,
		Status:              true,
	})
	if err != nil {
		return uuid.Nil, serviceError(err, codes.Internal)
	}

	return id, nil
}

func (s *categoryService) GetList(options dto.FilterOptions) ([]models.Category, int, error) {
	categories, total, err := s.repo.Category.GetList(options)
	if err != nil {
		return nil, 0, serviceError(err, codes.Internal)
	}

	resultCategories := make([]models.Category, len(categories))
	for i := range categories {
		resultCategories[i] = dtoToModelCategory(categories[i])
	}

	return resultCategories, total, nil
}

func dtoToModelCategory(category dto.Category) models.Category {
	return models.Category{
		ID:                  category.ID,
		ParentID:            category.ParentID,
		Name:                category.Name,
		Type:                category.Type,
		Photo:               category.Photo,
		IsTop:               category.IsTop,
		ApplicationAreasImg: category.ApplicationAreasImg,
		Status:              category.Status,
	}
}
