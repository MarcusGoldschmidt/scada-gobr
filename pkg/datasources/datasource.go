package datasources

import (
	"context"
	"gorm.io/gorm"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/server"
	"scadagobr/pkg/shared"
)

type Datapoint interface {
	Id() shared.CommonId
	Name() string
}

type DataSource interface {
	Id() shared.CommonId
	Name() string
	IsEnable() bool
	GetDataPoints() []Datapoint
	CreateRuntime(ctx context.Context, persistence persistence.DataPointPersistence) (DataSourceRuntime, error)
}

type DataSourceRuntime interface {
	Run(ctx context.Context, shutdownCompleteChan chan shared.CommonId) error
	GetDataSource() DataSource
}

type DataSourceCommon struct {
	id       shared.CommonId
	name     string
	isEnable bool
}

func NewDataSourceCommon(id shared.CommonId, name string, isEnable bool) *DataSourceCommon {
	return &DataSourceCommon{id: id, name: name, isEnable: isEnable}
}

type DataPointCommon struct {
	id       shared.CommonId
	name     string
	isEnable bool
}

func NewDataPointCommon(id shared.CommonId, name string, isEnable bool) *DataPointCommon {
	return &DataPointCommon{id: id, name: name, isEnable: isEnable}
}

func LoadDataSources(db *gorm.DB, router *server.Router) []DataSource {
	return []DataSource{}
}
