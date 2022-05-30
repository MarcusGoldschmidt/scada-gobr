package runtime

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"sync"
)

type RuntimeManagerOptions struct {
	MaxRuntimeRetry int
}

type RuntimeManager struct {
	Logger      logger.Logger
	mutex       sync.RWMutex
	dataSources map[shared.CommonId]datasources.DataSourceRuntimeManager
	options     RuntimeManagerOptions

	persistence  persistence.DataPointPersistence
	timeProvider providers.TimeProvider
}

func NewRuntimeManager(logger logger.Logger, persistence persistence.DataPointPersistence) *RuntimeManager {
	return &RuntimeManager{
		Logger:       logger,
		mutex:        sync.RWMutex{},
		dataSources:  make(map[shared.CommonId]datasources.DataSourceRuntimeManager),
		persistence:  persistence,
		timeProvider: providers.UtcTimeProvider{},
		options:      RuntimeManagerOptions{MaxRuntimeRetry: 5},
	}
}

func (r *RuntimeManager) WithTimeProvider(provider providers.TimeProvider) {
	r.timeProvider = provider
}

func (r *RuntimeManager) WithOptions(opt RuntimeManagerOptions) {
	r.options = opt
}

func (r *RuntimeManager) AddDataSource(sources ...datasources.DataSourceRuntimeManager) {
	for _, source := range sources {
		r.dataSources[source.Id()] = source
	}
}

func (r *RuntimeManager) RemoveDataSource(id shared.CommonId) {
	delete(r.dataSources, id)
}

func (r *RuntimeManager) Run(ctx context.Context, id shared.CommonId) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	dataSource := r.dataSources[id]

	err := dataSource.Run(ctx)
	if err != nil {
		r.Logger.Errorf("datasource runtime %s stopped with error: %s", dataSource.Name(), err.Error())
		return
	}
}

func (r *RuntimeManager) RunAll(ctx context.Context) {
	for id, _ := range r.dataSources {
		r.Run(ctx, id)
	}
}

func (r *RuntimeManager) UpdateDataSource(ctx context.Context, ds datasources.DataSourceRuntimeManager) {
	_ = r.StopDataSource(ctx, ds.Id())
	r.RemoveDataSource(ds.Id())
	r.AddDataSource(ds)
	r.Run(ctx, ds.Id())
}

func (r *RuntimeManager) RestartDataSource(ctx context.Context, id shared.CommonId) error {
	r.Logger.Infof("Restarting datasource %s", id.String())
	err := r.StopDataSource(ctx, id)
	if err != nil {
		return err
	}
	r.Run(ctx, id)
	return nil
}

func (r *RuntimeManager) StopDataSource(ctx context.Context, id shared.CommonId) error {
	r.Logger.Infof("Shutdown datapoint runtime %s", id.String())

	datasourceManager, ok := r.dataSources[id]

	if !ok {
		return errors.New("datasource not found")
	}

	err := datasourceManager.Stop(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *RuntimeManager) StopAll(ctx context.Context) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(len(r.dataSources))

	for id := range r.dataSources {
		id := id
		go func() {
			defer wg.Done()

			err := r.StopDataSource(ctx, id)
			if err != nil {
				r.Logger.Warningf("error on stopping data source Id: %s, %s", id.String(), err)
				return
			}
		}()
	}

	wg.Wait()
}
