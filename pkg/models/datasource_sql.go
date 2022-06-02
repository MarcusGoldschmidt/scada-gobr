package models

import "time"

type DataSourceTypeSql struct {
	Period           time.Duration
	Driver           string
	Query            string
	ConnectionString string
}

type DataPointTypeSql struct {
	RowIdentifier string
}
