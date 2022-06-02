package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type DataSeries struct {
	TimeStamp   time.Time `gorm:"type:timestamp"`
	Value       float64
	DataPointId shared.CommonId
}
