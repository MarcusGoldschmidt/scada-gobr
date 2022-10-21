package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type DataSource struct {
	Id         shared.CommonId `json:"id" gorm:"type:uuid"`
	Name       string          `json:"name"`
	Data       []byte          `json:"-" gorm:"type:jsonb"`
	Type       DataSourceType  `json:"type"`
	DataPoints []*DataPoint    `json:"dataPoints" gorm:"foreignKey:DataSourceId" `
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
	TypeData   map[string]any  `json:"data" gorm:"-"`
}
