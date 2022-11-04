package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type DataSource struct {
	Id          shared.CommonId `json:"id" gorm:"type:uuid"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Data        []byte          `json:"-" gorm:"type:jsonb"`
	Type        DataSourceType  `json:"type"`
	DataPoints  []*DataPoint    `json:"dataPoints" gorm:"foreignKey:DataSourceId" `
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	TypeData    map[string]any  `json:"data" gorm:"-"`
}

func (ds *DataSource) FilterAvailableDataPoints() {
	filteredDataPoints := make([]*DataPoint, 0)

	for _, point := range ds.DataPoints {
		if point.IsEnable {
			filteredDataPoints = append(filteredDataPoints, point)
		}
	}

	ds.DataPoints = filteredDataPoints
}
