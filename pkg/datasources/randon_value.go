package datasources

import (
	"context"
	"math/rand"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/providers"
	"scadagobr/pkg/shared"
	"time"
)

type RandonValueDataPoint struct {
	*DataPointCommon
	InitialValue int64
	EndValue     int64
}

func NewRandonValueDataPoint(dataPointCommon *DataPointCommon, initialValue int64, endValue int64) *RandonValueDataPoint {
	return &RandonValueDataPoint{DataPointCommon: dataPointCommon, InitialValue: initialValue, EndValue: endValue}
}

func (r RandonValueDataPoint) Id() shared.CommonId {
	return r.id
}

func (r RandonValueDataPoint) Name() string {
	return r.name
}

type RandomValueWorker struct {
	period       time.Duration
	dataPoints   []*RandonValueDataPoint
	persistence  persistence.DataPointPersistence
	dataSourceId shared.CommonId
}

func (c *RandomValueWorker) DataSourceId() shared.CommonId {
	return c.dataSourceId
}

func NewRandomValueWorker(dataSourceId shared.CommonId, period time.Duration, dataPoints []*RandonValueDataPoint, persistence persistence.DataPointPersistence) *RandomValueWorker {
	return &RandomValueWorker{dataSourceId: dataSourceId, period: period, dataPoints: dataPoints, persistence: persistence}
}

func (c *RandomValueWorker) Work(ctx context.Context, errorChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.period):
			for _, point := range c.dataPoints {
				value := point.InitialValue + rand.Int63n(point.EndValue-point.InitialValue)

				currentTime := providers.DefaultTimeProvider.GetCurrentTime()

				series := shared.NewSeries(float64(value), currentTime)

				err := c.persistence.AddDataPointValue(ctx, c.dataSourceId, series)
				if err != nil {
					errorChan <- err
					return
				}
			}
		}
	}
}
