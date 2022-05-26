package persistence

import (
	"context"
	"errors"
	"scadagobr/pkg/shared"
)

type InMemoryPersistence struct {
	data map[shared.CommonId][]*shared.Series
}

func (f InMemoryPersistence) GetPointValues(id shared.CommonId) ([]*shared.Series, error) {
	if data, ok := f.data[id]; ok {
		return data, nil
	}

	return nil, errors.New("data source not found")
}

func (f InMemoryPersistence) AddDataPointValues(ctx context.Context, values []*shared.IdSeries) error {
	for _, value := range values {
		err := f.AddDataPointValue(ctx, value.Id, shared.NewSeries(value.Value, value.Timestamp))
		if err != nil {
			return err
		}
	}
	return nil
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
