package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type DataSeries struct {
	Timestamp   time.Time       `json:"timestamp" gorm:"type:timestamp"`
	Value       float64         `json:"value"`
	DataPointId shared.CommonId `json:"-"`
}

func NewDataSeries(timeStamp time.Time, value float64, dataPointId shared.CommonId) *DataSeries {
	return &DataSeries{Timestamp: timeStamp, Value: value, DataPointId: dataPointId}
}
