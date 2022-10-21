package gorm

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"gorm.io/gorm"
	"time"
)

type DataPointPersistenceGormImpl struct {
	db                  *gorm.DB
	hubManager          events.HubManager
	dataPointGroupCache map[shared.CommonId]string
}

func NewDataPointPersistenceGormImpl(db *gorm.DB, hubManager events.HubManager) *DataPointPersistenceGormImpl {
	return &DataPointPersistenceGormImpl{
		db:                  db,
		hubManager:          hubManager,
		dataPointGroupCache: map[shared.CommonId]string{},
	}
}

func (d *DataPointPersistenceGormImpl) GetGroupNameByDataPointId(ctx context.Context, datapointId shared.CommonId) (string, error) {
	db := d.db.WithContext(ctx)

	if value, ok := d.dataPointGroupCache[datapointId]; ok {
		return value, nil
	}

	query := `
SELECT 
	CONCAT(data_sources.name, '-', data_points.name) as Group
FROM data_series
INNER JOIN data_points ON data_points.id = data_series.data_point_id
INNER JOIN data_sources ON data_sources.id = data_points.data_source_id
WHERE data_points.id = ?
`

	var group string
	response := db.Raw(query, datapointId).Scan(&group)
	if response.Error != nil {
		return "", response.Error
	}

	d.dataPointGroupCache[datapointId] = group

	return group, response.Error
}

func (d DataPointPersistenceGormImpl) GetPointValuesByIds(ctx context.Context, id []shared.CommonId, begin time.Time, end time.Time) ([]*persistence.SeriesGroupIdentifier, error) {
	db := d.db.WithContext(ctx)

	query := `
SELECT 
	CONCAT(data_sources.name, '-', data_points.name) as Group
, 	timestamp AS Timestamp
, 	value AS Value 
FROM data_series
INNER JOIN data_points ON data_points.id = data_series.data_point_id
INNER JOIN data_sources ON data_sources.id = data_points.data_source_id
WHERE data_point_id IN (?) AND timestamp > ? AND timestamp < ?
`

	var values []*persistence.SeriesGroupIdentifier
	response := db.Raw(query, id, begin, end).Scan(&values)

	return values, response.Error
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

	d.sendNotificationSeriesCreated(ctx, dataSeries)

	return nil
}

func (d DataPointPersistenceGormImpl) AddDataPointValues(ctx context.Context, values []*models.DataSeries) error {
	db := d.db.WithContext(ctx)
	err := db.Create(&values).Error

	if err != nil {
		return err
	}

	d.sendNotificationSeriesCreated(ctx, values...)

	return nil
}

func (d DataPointPersistenceGormImpl) GetPointValues(ctx context.Context, id shared.CommonId, begin time.Time, end time.Time) ([]*shared.Series, error) {
	db := d.db.WithContext(ctx)

	var values []*shared.Series

	db.Model(&models.DataSeries{}).Where("time < ? AND ? < time ", begin, end).Scan(&values)

	return values, nil
}

func (d *DataPointPersistenceGormImpl) sendNotificationSeriesCreated(ctx context.Context, series ...*models.DataSeries) {
	for _, data := range series {
		groupName, err := d.GetGroupNameByDataPointId(ctx, data.DataPointId)
		if err != nil {
			continue
		}

		event := persistence.NewSeriesGroupIdentifier(data.Timestamp, data.Value, groupName)

		d.hubManager.SendMessage(events.DataSeriesInserter+data.DataPointId.String(), event)
	}
}
