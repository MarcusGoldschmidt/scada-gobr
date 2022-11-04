package purge

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"sync"
	"time"
)

type Manager struct {
	dataPointPersistence  persistence.DataPointPersistence
	dataSourcePersistence persistence.DataSourcePersistence
	timeProvider          providers.TimeProvider
	logger                logger.Logger
	intervalCheck         time.Duration
}

func NewManager(
	dataPointPersistence persistence.DataPointPersistence,
	dataSourcePersistence persistence.DataSourcePersistence,
	timeProvider providers.TimeProvider,
	logger logger.Logger,
	intervalCheck time.Duration,
) *Manager {
	return &Manager{
		dataPointPersistence:  dataPointPersistence,
		dataSourcePersistence: dataSourcePersistence,
		timeProvider:          timeProvider,
		logger:                logger,
		intervalCheck:         intervalCheck,
	}
}

func (m *Manager) Purge(ctx context.Context) error {
	points, err := m.dataPointPersistence.GetAllDataPoints(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, point := range points {
		if point.PurgeAfter == nil {
			continue
		}

		wg.Add(1)

		now := m.timeProvider.GetCurrentTime()

		go func(point *models.DataPoint) {
			defer wg.Done()

			purgeDate := now.Add(-(*point.PurgeAfter))

			err := m.dataPointPersistence.DeleteDataPointValueByPeriod(ctx, point.Id, time.Time{}, purgeDate)
			if err != nil {
				m.logger.Errorf("error on purge %s", err.Error())
			}
		}(point)
	}

	wg.Wait()

	return nil
}

func (m *Manager) Work(ctx context.Context) {
	ticker := time.NewTicker(m.intervalCheck)
	for {
		select {
		case <-ctx.Done():
			m.logger.Infof("%s", ctx.Err())
			return
		case <-ticker.C:
			err := m.Purge(ctx)
			if err != nil {
				m.logger.Errorf("%s", err.Error())
			}
		}
	}
}
