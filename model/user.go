package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	PhoneNumber  string
	FullName     string
	PasswordHash string
}
