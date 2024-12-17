package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID    `json:"id"`
	FullName    string       `json:"full_name"`
	PhoneNumber string       `json:"phone_number"`
	RoleID      uuid.UUID    `json:"role_id"`
	Role        Role         `json:"role"`
	Email       string       `json:"email"`
	Password    string       `json:"-"`
	Company     *UserCompany `json:"company"`
	CreatedAt   time.Time    `json:"created_at"`
}

type UserCompany struct {
	INN  string `json:"inn"`
	Name string `json:"name"`
	File string `json:"file"`
}

type CreateUser struct {
	FullName    string       `json:"full_name"`
	PhoneNumber string       `json:"phone_number"`
	RoleID      uuid.UUID    `json:"role_id"`
	Email       string       `json:"email"`
	Password    string       `json:"password"`
	Company     *UserCompany `json:"company"`
}
