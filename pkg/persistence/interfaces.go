package persistence

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"time"
)

type DataPointPersistence interface {
	AddDataPointValue(ctx context.Context, id shared.CommonId, value *shared.Series) error
	AddDataPointValues(ctx context.Context, values []*shared.IdSeries) error
	GetPointValues(ctx context.Context, id shared.CommonId, begin time.Time, end time.Time) ([]*shared.Series, error)
}

type DataSourcePersistence interface {
	GetDadaSourceById(ctx context.Context, id shared.CommonId) (*models.DataSource, error)
	GetDadaSources(ctx context.Context) ([]*models.DataSource, error)
}

type UserPersistence interface {
	GetUserById(context.Context, uuid.UUID) (*models.User, error)
	GetUserByEmail(context.Context, string) (*models.User, error)
	GetUserByName(context.Context, string) (*models.User, error)
	IsValidUsernameForUser(context.Context, string, uuid.UUID) (bool, error)
	GetUsers(context.Context, *shared.PaginationRequest) ([]*models.User, error)

	CreateUser(context.Context, *models.User) error
	UpdateUser(context.Context, *models.User) error
	DeleteUser(context.Context, uuid.UUID) error
}
