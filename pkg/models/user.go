package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID            uuid.UUID
	Name          string
	Email         *string
	HomeUrl       string
	Administrator bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
