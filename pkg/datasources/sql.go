package datasources

import (
	"context"
	"database/sql"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
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

type SqlWorker struct {
	Period     time.Duration
	DataPoints []*SqlDataPoint

	// Dabatase
	Driver           string
	Query            string
	ConnectionString string

	dataSourceId shared.CommonId

	Persistence persistence.DataPointPersistence
}

func NewSqlWorker(period time.Duration, dataPoints []*SqlDataPoint, driver string, query string, connectionString string, dataSourceId shared.CommonId, persistence persistence.DataPointPersistence) *SqlWorker {
	return &SqlWorker{Period: period, DataPoints: dataPoints, Driver: driver, Query: query, ConnectionString: connectionString, dataSourceId: dataSourceId, Persistence: persistence}
}

func (c *SqlWorker) DataSourceId() shared.CommonId {
	return c.dataSourceId
}

func (c *SqlWorker) Run(ctx context.Context, errorChan chan<- error) {
	db, err := sql.Open(c.Driver, c.ConnectionString)

	if err != nil {
		errorChan <- err
		return
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		errorChan <- err
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.Period):
			series, err := c.QueryDatabase(ctx, db)
			if err != nil {
				errorChan <- err
				return
			}
			err = c.Persistence.AddDataPointValues(ctx, series)
			if err != nil {
				errorChan <- err
				return
			}
		}
	}
}

func (c *SqlWorker) QueryDatabase(ctx context.Context, db *sql.DB) ([]*shared.IdSeries, error) {

	dict := make(map[string]*shared.CommonId)

	for _, point := range c.DataPoints {
		dict[point.rowIdentifier] = &point.id
	}

	rows, err := db.QueryContext(ctx, c.Query)

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
