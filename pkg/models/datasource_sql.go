package models

import "time"

type DataSourceTypeSql struct {
	Period           time.Duration `json:"period" validate:"required"`
	Driver           string        `json:"driver" validate:"required"`
	Query            string        `json:"query" validate:"required"`
	ConnectionString string        `json:"connectionString" validate:"required"`
}

type DataPointTypeSql struct {
	RowIdentifier string `json:"row_identifier" validate:"required"`
}
