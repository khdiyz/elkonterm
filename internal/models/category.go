package models

import (
	"github.com/google/uuid"
)

type Category struct {
	ID                  uuid.UUID       `json:"id"`
	ParentID            *uuid.UUID      `json:"parent_id"`
	Name                NameTranslation `json:"name"`
	Type                string          `json:"type"`
	Photo               string          `json:"photo"`
	IsTop               bool            `json:"is_top"`
	ApplicationAreasImg string          `json:"application_areas_img"`
	Status              bool            `json:"status"`
}

type CreateCategory struct {
	ParentID            string          `json:"parent_id"`
	Name                NameTranslation `json:"name" binding:"required"`
	Type                string          `json:"type" binding:"required"`
	Photo               string          `json:"photo"`
	IsTop               bool            `json:"is_top"`
	ApplicationAreasImg string          `json:"application_areas_img"`
}
