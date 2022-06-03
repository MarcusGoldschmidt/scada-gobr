package purge

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
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

	for _, point := range points {

		if point.PurgeAfter == nil {
			continue
		}

		now := m.timeProvider.GetCurrentTime()

		begin := now.Add(-(*point.PurgeAfter))

		err := m.dataPointPersistence.DeleteDataPointValueByPeriod(ctx, point.Id, begin, now)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Work(ctx context.Context) {
	ticker := time.NewTicker(m.intervalCheck)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := m.Purge(ctx)
			if err != nil {
				m.logger.Errorf("%s", err.Error())
				return
			}
		}
	}
}
