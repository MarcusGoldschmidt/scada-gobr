package gorm

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"gorm.io/gorm"
	"time"
)

type DataPointPersistenceGormImpl struct {
	db *gorm.DB
}

func (d DataPointPersistenceGormImpl) AddDataPointValue(ctx context.Context, id shared.CommonId, value *shared.Series) error {
	db := d.db.WithContext(ctx)
	return db.Create(value).Error
}

func (d DataPointPersistenceGormImpl) AddDataPointValues(ctx context.Context, values []*models.DataSeries) error {
	db := d.db.WithContext(ctx)
	return db.Create(&values).Error
}

func (d DataPointPersistenceGormImpl) GetPointValues(ctx context.Context, id shared.CommonId, begin time.Time, end time.Time) ([]*shared.Series, error) {
	db := d.db.WithContext(ctx)

	var values []*shared.Series

	db.Model(&models.DataSeries{}).Where("time < ? AND ? < time ", begin, end).Scan(&values)

	return values, nil
}
