package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type DataSource struct {
	ID   shared.CommonId `gorm:"type:uuid"`
	Name string
	Data []byte `gorm:"type:jsonb"`

	Type DataSourceType

	DataPoints []*DataPoint `gorm:"foreignKey:DataSourceId"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
