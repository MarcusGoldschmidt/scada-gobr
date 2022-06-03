package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type DataSeries struct {
	Timestamp   time.Time `gorm:"type:timestamp"`
	Value       float64
	DataPointId shared.CommonId
}

func NewDataSeries(timeStamp time.Time, value float64, dataPointId shared.CommonId) *DataSeries {
	return &DataSeries{Timestamp: timeStamp, Value: value, DataPointId: dataPointId}
}
