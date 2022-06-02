package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid"`
	PasswordHash  string
	Name          string `gorm:"unique"`
	Email         *string
	HomeUrl       string
	Administrator bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
