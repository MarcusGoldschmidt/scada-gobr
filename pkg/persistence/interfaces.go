package persistence

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"time"
)

type SeriesGroupIdentifier struct {
	Timestamp time.Time `json:"timestamp" gorm:"type:timestamp"`
	Value     float64   `json:"value"`
	Group     string    `json:"group"`
}

func NewSeriesGroupIdentifier(timestamp time.Time, value float64, group string) *SeriesGroupIdentifier {
	return &SeriesGroupIdentifier{Timestamp: timestamp, Value: value, Group: group}
}

type DataPointPersistence interface {
	AddDataPointValue(ctx context.Context, id shared.CommonId, value *shared.Series) error
	AddDataPointValues(ctx context.Context, values []*models.DataSeries) error
	DeleteDataPointValueByPeriod(ctx context.Context, id shared.CommonId, begin time.Time, end time.Time) error

	GetPointValues(ctx context.Context, id shared.CommonId, begin time.Time, end time.Time) ([]*shared.Series, error)
	GetPointValuesByIds(ctx context.Context, id []shared.CommonId, begin time.Time, end time.Time) ([]*SeriesGroupIdentifier, error)
	GetGroupNameByDataPointId(ctx context.Context, datapointId shared.CommonId) (string, error)

	CreateDataPoint(ctx context.Context, dataPoint *models.DataPoint) error
	GetDataPointById(ctx context.Context, id shared.CommonId) (*models.DataPoint, error)
	GetAllDataPoints(ctx context.Context) ([]*models.DataPoint, error)
	UpdateDataPoint(ctx context.Context, dataPoint *models.DataPoint) error
	DeleteDataPoint(ctx context.Context, dataSourceId shared.CommonId, dataPointId shared.CommonId) error
}

type DataSourcePersistence interface {
	GetDataSourceById(ctx context.Context, id shared.CommonId) (*models.DataSource, error)
	GetDataSources(ctx context.Context) ([]*models.DataSource, error)

	GetDataSourceByName(ctx context.Context, name string) (*models.DataSource, error)

	GetDataPoints(ctx context.Context, id shared.CommonId) ([]*models.DataPoint, error)

	CreateDataSource(ctx context.Context, dataSource *models.DataSource) error
	DeleteDataSource(ctx context.Context, id shared.CommonId) error
}

type UserPersistence interface {
	GetUserById(context.Context, uuid.UUID) (*models.User, error)
	GetUserByEmail(context.Context, string) (*models.User, error)
	GetUserByName(context.Context, string) (*models.User, error)
	IsValidUsernameForUser(context.Context, string, uuid.UUID) (bool, error)
	GetUsers(context.Context, *shared.PaginationRequest) ([]*models.User, error)

	CreateAdminUser(ctx context.Context, email string, password string) error
	CreateUser(context.Context, *models.User) error
	UpdateUser(context.Context, *models.User) error
	DeleteUser(context.Context, uuid.UUID) error
}

type ViewPersistence interface {
	GetViewById(context.Context, uuid.UUID) (*models.View, error)
	GetAllViews(context.Context) ([]*models.View, error)

	GetViewComponentById(context.Context, uuid.UUID) (*models.ViewComponent, error)

	AttachViewComponents(context.Context, ...*models.ViewComponent) error
	DeleteViewComponent(ctx context.Context, viewId uuid.UUID, componentId uuid.UUID) error

	CreateView(context.Context, *models.View) error
	UpdateView(context.Context, *models.View) error
	DeleteView(context.Context, uuid.UUID) error
}
