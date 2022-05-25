package datasources

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/shared"
	"time"
)

type SqlDataPoint struct {
	*DataPointCommon
	rowIdentifier string
}

func (r SqlDataPoint) Id() shared.CommonId {
	return r.id
}

func (r SqlDataPoint) Name() string {
	return r.name
}

type SqlDataSource struct {
	*DataSourceCommon
	period     time.Duration
	dataPoints []*SqlDataPoint

	// Dabatase
	driver           string
	query            string
	connectionString string
}

type SqlDataSourceRuntime struct {
	id          shared.CommonId
	dataSource  SqlDataSource
	persistence persistence.DataPointPersistence
}

func (r *SqlDataSource) AddDataPoint(dp *SqlDataPoint) {
	r.dataPoints = append(r.dataPoints, dp)
}

func (r SqlDataSource) Id() shared.CommonId {
	return r.id
}

func (r SqlDataSource) Name() string {
	return r.name
}

func (r SqlDataSource) IsEnable() bool {
	return r.isEnable
}

func (r SqlDataSource) GetDataPoints() []Datapoint {
	datapoint := make([]Datapoint, len(r.dataPoints))
	for i, v := range r.dataPoints {
		datapoint[i] = Datapoint(v)
	}
	return datapoint
}

func (r SqlDataSource) CreateRuntime(ctx context.Context, p persistence.DataPointPersistence) (DataSourceRuntime, error) {
	rt := SqlDataSourceRuntime{uuid.New(), r, p}
	return rt, nil
}

func (c SqlDataSourceRuntime) GetDataSource() DataSource {
	return c.dataSource
}

func (c SqlDataSourceRuntime) Run(ctx context.Context, shutdownCompleteChan chan shared.CommonId) error {
	defer func() {
		shutdownCompleteChan <- c.id
	}()

	db, err := sql.Open(c.dataSource.driver, c.dataSource.connectionString)

	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	defer db.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(c.dataSource.period):
			series, err := c.QueryDatabase(ctx, db)
			if err != nil {
				return err
			}
			err = c.persistence.AddDataPointValues(ctx, series)
			if err != nil {
				return err
			}
		}
	}
}

func (c SqlDataSourceRuntime) QueryDatabase(ctx context.Context, db *sql.DB) ([]*shared.IdSeries, error) {

	dict := make(map[string]*shared.CommonId)

	for _, point := range c.dataSource.dataPoints {
		dict[point.rowIdentifier] = &point.id
	}

	rows, err := db.QueryContext(ctx, c.dataSource.query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*shared.IdSeries, 0)

	for rows.Next() {
		var identifier string
		var value float64
		var timestamp time.Time

		if err := rows.Scan(&identifier, value, timestamp); err != nil {
			return nil, err
		}

		if id, ok := dict[identifier]; ok {
			result = append(result, shared.NewIdSeries(*id, value, timestamp))
		}
	}

	return result, nil
}
