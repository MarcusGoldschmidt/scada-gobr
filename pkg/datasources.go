package pkg

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"gorm.io/gorm"
)

func LoadDataSourceRuntimeManager(db *gorm.DB, router *server.Router) []datasources.DataSourceRuntimeManager {
	return []datasources.DataSourceRuntimeManager{}
}
