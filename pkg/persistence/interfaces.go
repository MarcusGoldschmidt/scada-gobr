package persistence

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
)

type DataPointPersistence interface {
	AddDataPointValue(ctx context.Context, id shared.CommonId, value *shared.Series) error
	AddDataPointValues(ctx context.Context, values []*shared.IdSeries) error
	GetPointValues(id shared.CommonId) ([]*shared.Series, error)
}

type UserPersistence interface {
	GetUserById(uuid.UUID) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	GetUserByName(string) (*models.User, error)
}
