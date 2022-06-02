package in_memory

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"time"
)

type InMemoryPersistence struct {
	data map[shared.CommonId][]*shared.Series
}

func (f InMemoryPersistence) AddDataPointValues(ctx context.Context, values []*shared.IdSeries) error {
	for _, value := range values {
		err := f.AddDataPointValue(ctx, value.Id, &shared.Series{
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
			result = append(result, value)
		}
	}

	return result, nil
}

func NewInMemoryPersistence() *InMemoryPersistence {
	return &InMemoryPersistence{data: make(map[shared.CommonId][]*shared.Series)}
}

func (f InMemoryPersistence) AddDataPointValue(ctx context.Context, id shared.CommonId, value *shared.Series) error {
	if f.data[id] == nil {
		f.data[id] = []*shared.Series{value}
	} else {
		f.data[id] = append(f.data[id], value)
	}
	return nil
}
