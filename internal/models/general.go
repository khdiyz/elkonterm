package models

import (
	"time"

	"github.com/google/uuid"
)

type NameTranslation struct {
	Uz string `json:"uz"`
	En string `json:"en"`
	Ru string `json:"ru"`
}

type Token struct {
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Pagination struct {
	Page       int `json:"page"  default:"1"`
	Limit      int `json:"limit" default:"10"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}

type IdResponse struct {
	ID uuid.UUID `json:"id"`
}
