package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type DataSourceTypeHttpRequest struct {
	Period           time.Duration     `json:"period"`
	BaseUrl          string            `json:"baseUrl"`
	Encoding         string            `json:"encoding"`
	Method           string            `json:"method"`
	Headers          map[string]string `json:"headers"`
	BodyTemplate     *string           `json:"bodyTemplate"`
	ForEachDataPoint bool              `json:"forEachDataPoint"`
	dataSourceId     shared.CommonId   `json:"dataSourceId"`
}

type DataPointTypeHttpRequest struct {
	RowIdentifier string
	DateFormat    string
}
