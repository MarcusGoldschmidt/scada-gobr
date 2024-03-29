package shared

import (
	"github.com/google/uuid"
	"time"
)

type CommonId = uuid.UUID

func NewCommonId() CommonId {
	return uuid.New()
}

func ParseCommonId(id string) (CommonId, error) {
	return uuid.Parse(id)
}

type Series struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func NewSeries(value float64, timestamp time.Time) *Series {
	return &Series{Value: value, Timestamp: timestamp}
}
