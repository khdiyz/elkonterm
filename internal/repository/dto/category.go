package dto

import (
	"elkonterm/internal/models"

	"github.com/google/uuid"
)

type Category struct {
	ID                  uuid.UUID
	ParentID            *uuid.UUID
	Name                models.NameTranslation
	Type                string
	Photo               string
	IsTop               bool
	ApplicationAreasImg string
	Status              bool
}

type CreateCategory struct {
	ParentID            *uuid.UUID
	Name                models.NameTranslation
	Type                string
	Photo               string
	IsTop               bool
	ApplicationAreasImg string
	Status              bool
}
