package models

import "time"

type DataSourceTypeRandomValue struct {
	Period time.Duration `json:"period"`
}

type DataPointTypeRandomValue struct {
	InitialValue int64 `json:"initialValue"`
	EndValue     int64 `json:"endValue"`
}
