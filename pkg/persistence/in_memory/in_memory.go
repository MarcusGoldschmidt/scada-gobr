package in_memory

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type InMemoryPersistence struct {
	series     map[shared.CommonId][]*models.DataSeries
	dataPoints map[shared.CommonId]*models.DataPoint
}

func (f *InMemoryPersistence) GetGroupNameByDataPointId(ctx context.Context, datapointId shared.CommonId) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (f *InMemoryPersistence) GetPointValuesByIds(ctx context.Context, ids []shared.CommonId, begin time.Time, end time.Time) ([]*persistence.SeriesGroupIdentifier, error) {
	var result []*persistence.SeriesGroupIdentifier

	for _, id := range ids {
		if points, ok := f.series[id]; ok {
			for _, point := range points {
				if point.Timestamp.After(begin) && point.Timestamp.Before(end) {
					temp := persistence.NewSeriesGroupIdentifier(point.Timestamp, point.Value, "")
					result = append(result, temp)
				}
			}
		}
	}

	return result, nil
}

func (f *InMemoryPersistence) DeleteDataPointValueByPeriod(ctx context.Context, id shared.CommonId, begin time.Time, end time.Time) error {
	if series, ok := f.series[id]; ok {
		var result []*models.DataSeries
		for _, dataSeries := range series {
			if !(dataSeries.Timestamp.After(begin) && dataSeries.Timestamp.Before(end)) {
				result = append(result, dataSeries)
			}
		}
		f.series[id] = result
	}

	return nil
}

func (f *InMemoryPersistence) CreateDataPoint(ctx context.Context, dataPoint *models.DataPoint) error {
	f.dataPoints[dataPoint.Id] = dataPoint
	return nil
}

func (f *InMemoryPersistence) GetDataPointById(ctx context.Context, id shared.CommonId) (*models.DataPoint, error) {
	if dp, ok := f.dataPoints[id]; ok {
		return dp, nil
	}

	return nil, errors.New("dataPoint not found")
}

func (f *InMemoryPersistence) GetAllDataPoints(ctx context.Context) ([]*models.DataPoint, error) {
	var dataPoints []*models.DataPoint

	for _, point := range f.dataPoints {
		dataPoints = append(dataPoints, point)
	}
	return dataPoints, nil
}

func (f InMemoryPersistence) UpdateDataPoint(ctx context.Context, dataPoint *models.DataPoint) error {
	f.dataPoints[dataPoint.Id] = dataPoint
	return nil
}

func (f InMemoryPersistence) DeleteDataPoint(ctx context.Context, dataSourceId shared.CommonId, dataPointId shared.CommonId) error {
	delete(f.dataPoints, dataPointId)
	return nil
}

func (f InMemoryPersistence) AddDataPointValues(ctx context.Context, values []*models.DataSeries) error {
	for _, value := range values {
		err := f.AddDataPointValue(ctx, value.DataPointId, &shared.Series{
			Value:     value.Value,
			Timestamp: value.Timestamp,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (f InMemoryPersistence) GetPointValues(ctx context.Context, id shared.CommonId, begin time.Time, end time.Time) ([]*shared.Series, error) {
	result := make([]*shared.Series, 0)

	//Filter by period
	for _, value := range f.series[id] {
		if value.Timestamp.After(begin) && value.Timestamp.Before(end) {
			result = append(result, shared.NewSeries(value.Value, value.Timestamp))
		}
	}

	return result, nil
}

func NewInMemoryPersistence() *InMemoryPersistence {
	return &InMemoryPersistence{series: make(map[shared.CommonId][]*models.DataSeries)}
}

func (f InMemoryPersistence) AddDataPointValue(ctx context.Context, id shared.CommonId, value *shared.Series) error {
	series := models.NewDataSeries(value.Timestamp, value.Value, id)

	if f.series[id] == nil {
		f.series[id] = []*models.DataSeries{series}
	} else {
		f.series[id] = append(f.series[id], series)
	}
	return nil
}
