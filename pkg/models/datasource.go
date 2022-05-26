package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type DataSource struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid"`
	Name      string
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
