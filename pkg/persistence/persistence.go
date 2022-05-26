package persistence

import (
	"context"
	"scadagobr/pkg/shared"
)

type DataPointPersistence interface {
	AddDataPointValue(ctx context.Context, id shared.CommonId, value *shared.Series) error
	AddDataPointValues(ctx context.Context, values []*shared.IdSeries) error

	GetPointValues(id shared.CommonId) ([]*shared.Series, error)
}
