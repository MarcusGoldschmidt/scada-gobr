package shared

import (
	"github.com/google/uuid"
	"time"
)

type CommonId = uuid.UUID

func NewCommonId() CommonId {
	return uuid.New()
}

type Series struct {
	Value     float64
	Timestamp time.Time
}

type IdSeries struct {
	Id        CommonId
	Value     float64
	Timestamp time.Time
}

func NewIdSeries(id CommonId, value float64, timestamp time.Time) *IdSeries {
	return &IdSeries{Id: id, Value: value, Timestamp: timestamp}
}

func NewSeries(value float64, timestamp time.Time) *Series {
	return &Series{Value: value, Timestamp: timestamp}
}
