package models

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
)

type DataSourceTypeHttpServer struct {
	Endpoint     string
	User         string
	PasswordHash string
	dataSourceId shared.CommonId
}

type DataPointTypeHttpServer struct {
	RowIdentifier string
	DateFormat    string
}
