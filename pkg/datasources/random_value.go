package datasources

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"math/rand"
	"time"
)

type RandomValueDataPoint struct {
	common       *DataPointCommon
	InitialValue int64
	EndValue     int64
}

func NewRandomValueDataPoint(dataPointCommon *DataPointCommon, initialValue int64, endValue int64) *RandomValueDataPoint {
	return &RandomValueDataPoint{common: dataPointCommon, InitialValue: initialValue, EndValue: endValue}
}

func (r RandomValueDataPoint) Id() shared.CommonId {
	return r.common.Id
}

func (r RandomValueDataPoint) Name() string {
	return r.common.Name
}

type RandomValueWorker struct {
	period       time.Duration
	dataPoints   []*RandomValueDataPoint
	persistence  persistence.DataPointPersistence
	dataSourceId shared.CommonId
}

func (c *RandomValueWorker) DataSourceId() shared.CommonId {
	return c.dataSourceId
}

func NewRandomValueWorker(dataSourceId shared.CommonId, period time.Duration, dataPoints []*RandomValueDataPoint, persistence persistence.DataPointPersistence) *RandomValueWorker {
	return &RandomValueWorker{dataSourceId: dataSourceId, period: period, dataPoints: dataPoints, persistence: persistence}
}

func (c *RandomValueWorker) Run(ctx context.Context, errorChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.period):
			for _, point := range c.dataPoints {
				value := point.InitialValue + rand.Int63n(point.EndValue-point.InitialValue)

				currentTime := providers.GetTimeProviderFromCtx(ctx).GetCurrentTime()

				series := shared.NewSeries(float64(value), currentTime)

				err := c.persistence.AddDataPointValue(ctx, point.Id(), series)
				if err != nil {
					errorChan <- err
					return
				}
			}
		}
	}
}
