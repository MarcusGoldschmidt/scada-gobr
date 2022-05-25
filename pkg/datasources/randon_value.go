package datasources

import (
	"context"
	"github.com/google/uuid"
	"math/rand"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/providers"
	"scadagobr/pkg/shared"
	"time"
)

type RandonValueDataPoint struct {
	*DataPointCommon
	initialValue int64
	endValue     int64
}

func NewRandonValueDataPoint(dataPointCommon *DataPointCommon, initialValue int64, endValue int64) *RandonValueDataPoint {
	return &RandonValueDataPoint{DataPointCommon: dataPointCommon, initialValue: initialValue, endValue: endValue}
}

func (r RandonValueDataPoint) Id() shared.CommonId {
	return r.id
}

func (r RandonValueDataPoint) Name() string {
	return r.name
}

type RandonValueDataSource struct {
	*DataSourceCommon
	period     time.Duration
	dataPoints []*RandonValueDataPoint
}

type RandonValueDataSourceRuntime struct {
	id          shared.CommonId
	dataSource  RandonValueDataSource
	persistence persistence.DataPointPersistence
}

func NewRandonValueDataSource(dataSourceCommon *DataSourceCommon, period time.Duration) *RandonValueDataSource {
	return &RandonValueDataSource{dataSourceCommon, period, []*RandonValueDataPoint{}}
}

func (r *RandonValueDataSource) AddDataPoint(dp *RandonValueDataPoint) {
	r.dataPoints = append(r.dataPoints, dp)
}

func (r RandonValueDataSource) Id() shared.CommonId {
	return r.id
}

func (r RandonValueDataSource) Name() string {
	return r.name
}

func (r RandonValueDataSource) IsEnable() bool {
	return r.isEnable
}

func (r RandonValueDataSource) GetDataPoints() []Datapoint {
	datapoint := make([]Datapoint, len(r.dataPoints))
	for i, v := range r.dataPoints {
		datapoint[i] = Datapoint(v)
	}
	return datapoint
}

func (r RandonValueDataSource) CreateRuntime(ctx context.Context, p persistence.DataPointPersistence) (DataSourceRuntime, error) {
	rt := RandonValueDataSourceRuntime{uuid.New(), r, p}
	return rt, nil
}

func (c RandonValueDataSourceRuntime) GetDataSource() DataSource {
	return c.dataSource
}

func (c RandonValueDataSourceRuntime) Run(ctx context.Context, shutdownCompleteChan chan shared.CommonId) error {
	defer func() {
		shutdownCompleteChan <- c.id
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(c.dataSource.period):
			for _, point := range c.dataSource.dataPoints {
				value := point.initialValue + rand.Int63n(point.endValue-point.initialValue)

				currentTime := providers.DefaultTimeProvider.GetCurrentTime()

				series := shared.NewSeries(float64(value), currentTime)

				err := c.persistence.AddDataPointValue(ctx, point.Id(), series)
				if err != nil {
					return err
				}
			}
		}
	}
}
