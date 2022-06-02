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
	Id shared.CommonId `gorm:"type:uuid"`

	DataSourceId shared.CommonId

	Name       string
	IsEnable   bool
	Unit       string
	PurgeAfter *time.Duration
	Type       DataSourceType
	Data       []byte `gorm:"type:jsonb"`
}
