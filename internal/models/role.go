package models

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Data []any     `json:"data,omitempty"`
}
