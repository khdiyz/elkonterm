package dto

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `db:"id"`
	FullName    string    `db:"full_name"`
	PhoneNumber string    `db:"phone_number"`
	RoleID      uuid.UUID `db:"role_id"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Company     []byte    `db:"company"`
	Status      bool      `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}

type CreateUser struct {
	FullName    string
	PhoneNumber string
	RoleID      uuid.UUID
	Email       string
	Password    string
	Company     []byte
	Status      bool
}
