package pkg

import (
	"gorm.io/gorm"
	"scadagobr/pkg/datasources"
	"scadagobr/pkg/server"
)

func LoadDataSourceRuntimeManager(db *gorm.DB, router *server.Router) []datasources.DataSourceRuntimeManager {
	return []datasources.DataSourceRuntimeManager{}
}
