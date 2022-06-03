package gorm

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"gorm.io/gorm"
)

type DataSourcePersistenceGormImpl struct {
	db *gorm.DB
}

func (d DataSourcePersistenceGormImpl) GetDataPoints(ctx context.Context, id shared.CommonId) ([]*models.DataPoint, error) {
	db := d.db.WithContext(ctx)

	var data []*models.DataPoint

	err := db.Model(&models.DataPoint{}).Where("data_source_id = ?", id).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (d DataSourcePersistenceGormImpl) CreateDataSource(ctx context.Context, dataSource *models.DataSource) error {
	db := d.db.WithContext(ctx)
	return db.Create(dataSource).Error
}

func (d DataSourcePersistenceGormImpl) DeleteDataSource(ctx context.Context, id shared.CommonId) error {
	db := d.db.WithContext(ctx)
	return db.Delete(&models.DataSource{ID: id}).Error
}

func NewDataSourcePersistenceGormImpl(db *gorm.DB) *DataSourcePersistenceGormImpl {
	return &DataSourcePersistenceGormImpl{db: db}
}

func (d DataSourcePersistenceGormImpl) GetDadaSourceById(ctx context.Context, id shared.CommonId) (*models.DataSource, error) {
	db := d.db.WithContext(ctx)
	return getById[models.DataSource](db, id)
}

func (d DataSourcePersistenceGormImpl) GetDadaSources(ctx context.Context) ([]*models.DataSource, error) {
	db := d.db.WithContext(ctx)

	var data []*models.DataSource
	err := db.Model(&models.DataSource{}).Preload("DataPoints").Find(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}
