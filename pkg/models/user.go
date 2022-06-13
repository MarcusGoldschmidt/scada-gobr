package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid"`
	PasswordHash  string    `json:"passwordHash"`
	Name          string    `json:"name" gorm:"unique"`
	Email         *string   `json:"email"`
	HomeUrl       string    `json:"homeUrl"`
	Administrator bool      `json:"administrator"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
