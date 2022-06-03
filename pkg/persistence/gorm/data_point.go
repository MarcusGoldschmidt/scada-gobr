package gorm

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"gorm.io/gorm"
	"time"
)

type DataPointPersistenceGormImpl struct {
	db         *gorm.DB
	hubManager events.HubManager
}

func NewDataPointPersistenceGormImpl(db *gorm.DB, hubManager events.HubManager) *DataPointPersistenceGormImpl {
	return &DataPointPersistenceGormImpl{db: db, hubManager: hubManager}
}

func (d DataPointPersistenceGormImpl) DeleteDataPointValueByPeriod(ctx context.Context, id shared.CommonId, begin time.Time, end time.Time) error {
	db := d.db.WithContext(ctx)
	return db.Delete(&models.DataSeries{}, "time > ? AND ? > time AND data_point_id = ?", begin, end, id).Error
}

func (d DataPointPersistenceGormImpl) GetAllDataPoints(ctx context.Context) ([]*models.DataPoint, error) {
	db := d.db.WithContext(ctx)
	return listAll[models.DataPoint](db)
}

func (d DataPointPersistenceGormImpl) CreateDataPoint(ctx context.Context, dataPoint *models.DataPoint) error {
	db := d.db.WithContext(ctx)
	return db.Create(dataPoint).Error
}

func (d DataPointPersistenceGormImpl) GetDataPointById(ctx context.Context, id shared.CommonId) (*models.DataPoint, error) {
	db := d.db.WithContext(ctx)
	return getById[models.DataPoint](db, id)
}

func (d DataPointPersistenceGormImpl) UpdateDataPoint(ctx context.Context, dataPoint *models.DataPoint) error {
	db := d.db.WithContext(ctx)
	return db.Save(dataPoint).Error
}

func (d DataPointPersistenceGormImpl) DeleteDataPoint(ctx context.Context, dataSourceId shared.CommonId, dataPointId shared.CommonId) error {
	db := d.db.WithContext(ctx)
	return db.Delete(&models.DataPoint{Id: dataPointId, DataSourceId: dataSourceId}).Error
}

func (d DataPointPersistenceGormImpl) AddDataPointValue(ctx context.Context, id shared.CommonId, value *shared.Series) error {
	db := d.db.WithContext(ctx)

	dataSeries := models.NewDataSeries(value.Timestamp, value.Value, id)

	err := db.Create(dataSeries).Error

	if err != nil {
		return err
	}

	d.sendNotificationSeriesCreated(dataSeries)

	return nil
}

func (d DataPointPersistenceGormImpl) AddDataPointValues(ctx context.Context, values []*models.DataSeries) error {
	db := d.db.WithContext(ctx)
	err := db.Create(&values).Error

	if err != nil {
		return err
	}

	d.sendNotificationSeriesCreated(values...)

	return nil
}

func (d DataPointPersistenceGormImpl) GetPointValues(ctx context.Context, id shared.CommonId, begin time.Time, end time.Time) ([]*shared.Series, error) {
	db := d.db.WithContext(ctx)

	var values []*shared.Series

	db.Model(&models.DataSeries{}).Where("time < ? AND ? < time ", begin, end).Scan(&values)

	return values, nil
}

func (d DataPointPersistenceGormImpl) sendNotificationSeriesCreated(series ...*models.DataSeries) {
	for _, data := range series {
		d.hubManager.SendMessage(events.DataSeriesInserter+data.DataPointId.String(), data)
	}
}
