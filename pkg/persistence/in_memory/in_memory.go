package in_memory

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type InMemoryPersistence struct {
	data map[shared.CommonId][]*models.DataSeries
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
	for _, value := range f.data[id] {
		if value.Timestamp.After(begin) && value.Timestamp.Before(end) {
			result = append(result, shared.NewSeries(value.Value, value.Timestamp))
		}
	}

	return result, nil
}

func NewInMemoryPersistence() *InMemoryPersistence {
	return &InMemoryPersistence{data: make(map[shared.CommonId][]*models.DataSeries)}
}

func (f InMemoryPersistence) AddDataPointValue(ctx context.Context, id shared.CommonId, value *shared.Series) error {
	series := models.NewDataSeries(value.Timestamp, value.Value, id)

	if f.data[id] == nil {
		f.data[id] = []*models.DataSeries{series}
	} else {
		f.data[id] = append(f.data[id], series)
	}
	return nil
}
