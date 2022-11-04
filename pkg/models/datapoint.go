package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type DataPointLoggingType int

const (
	AllData DataPointLoggingType = iota
	WhenValueChanges
	TimeStampChanges
)

type DataPoint struct {
	Id           shared.CommonId      `json:"id" gorm:"type:uuid"`
	DataSourceId shared.CommonId      `json:"dataSourceId"`
	Name         string               `json:"name"`
	Description  *string              `json:"description"`
	IsEnable     bool                 `json:"isEnable"`
	Unit         string               `json:"unit"`
	PurgeAfter   *time.Duration       `json:"purgeAfter"`
	Type         DataSourceType       `json:"type"`
	Data         []byte               `json:"-" gorm:"type:jsonb"`
	TypeData     map[string]any       `json:"data" gorm:"-"`
	LoggingType  DataPointLoggingType `json:"loggingType"`
	CreatedAt    time.Time            `json:"createdAt"`
	UpdatedAt    time.Time            `json:"updatedAt"`
}
