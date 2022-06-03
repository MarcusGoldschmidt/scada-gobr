package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
)

type DataSourceTypeHttpServer struct {
	Endpoint     string          `json:"endpoint"`
	User         string          `json:"user"`
	PasswordHash string          `json:"password_hash"`
	DataSourceId shared.CommonId `json:"data_source_id"`
}

type DataPointTypeHttpServer struct {
	RowIdentifier string `json:"row_identifier"`
	DateFormat    string `json:"date_format"`
}
