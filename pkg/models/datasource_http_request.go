package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type DataSourceTypeHttpRequest struct {
	Period time.Duration

	BaseUrl      string
	Encoding     string
	Method       string
	Headers      map[string]string
	BodyTemplate *string

	ForEachDataPoint bool

	dataSourceId shared.CommonId
}

type DataPointTypeHttpRequest struct {
	RowIdentifier string
	DateFormat    string
}
