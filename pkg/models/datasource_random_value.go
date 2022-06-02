package models

import "time"

type DataSourceTypeRandomValue struct {
	Period time.Duration
}

type DataPointTypeRandomValue struct {
	InitialValue int64
	EndValue     int64
}
