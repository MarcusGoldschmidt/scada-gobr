package runtime

import (
	"context"
	"scadagobr/pkg/datasources"
	"scadagobr/pkg/logger"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/providers"
	"scadagobr/pkg/shared"
	"sync"
	"time"
)

type RuntimeManagerOptions struct {
	MaxRuntimeRetry int
}

type dataSourcesRuntimeControl struct {
	sync.Mutex
	shutdown        func()
	confirmShutdown chan shared.CommonId
}

type RuntimeManager struct {
	Logger                    logger.Logger
	mutex                     sync.RWMutex
	dataSources               map[shared.CommonId]datasources.DataSource
	dataSourcesRuntimeControl map[shared.CommonId]*dataSourcesRuntimeControl
	options                   RuntimeManagerOptions

	persistence  persistence.DataPointPersistence
	timeProvider providers.TimeProvider
}

func NewRuntimeManager(logger logger.Logger, persistence persistence.DataPointPersistence) *RuntimeManager {
	return &RuntimeManager{
		Logger:                    logger,
		mutex:                     sync.RWMutex{},
		dataSources:               make(map[shared.CommonId]datasources.DataSource),
		dataSourcesRuntimeControl: make(map[shared.CommonId]*dataSourcesRuntimeControl),
		persistence:               persistence,
		timeProvider:              providers.UtcTimeProvider{},
		options:                   RuntimeManagerOptions{MaxRuntimeRetry: 5},
	}
}

func (r *RuntimeManager) WithTimeProvider(provider providers.TimeProvider) {
	r.timeProvider = provider
}

func (r *RuntimeManager) WithOptions(opt RuntimeManagerOptions) {
	r.options = opt
}

func (r *RuntimeManager) AddDataSource(sources ...datasources.DataSource) {
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

	if !dataSource.IsEnable() {
		return
	}

	runtime, err := dataSource.CreateRuntime(ctx, r.persistence)

	if err != nil {
		r.Logger.Errorf("Error to create data source runtime %s reason: %s", dataSource.Id(), err.Error())
		return
	}

	shutdown := make(chan shared.CommonId)
	go func() {
		for errCount := 0; errCount < r.options.MaxRuntimeRetry; errCount++ {
			ctx, cancel := context.WithCancel(ctx)
			r.createShutdownManager(dataSource.Id(), cancel, shutdown)

			r.Logger.Infof("Starting datasource runtime %s", dataSource.Name())

			err := runtime.Run(ctx, shutdown)
			if err != nil {
				r.Logger.Warningf("Restarting datapoint runtime %s with error %s", dataSource.Name(), err.Error())
			}
			r.StopDataSource(dataSource.Id())
		}

		r.Logger.Warningf("datasource runtime %s stopped with max restart option", dataSource.Name())
	}()
}

func (r *RuntimeManager) RunAll(ctx context.Context) {
	for id, _ := range r.dataSources {
		r.Run(ctx, id)
	}
}

func (r *RuntimeManager) createShutdownManager(id shared.CommonId, cancel func(), shutdown chan shared.CommonId) {
	r.dataSourcesRuntimeControl[id] = &dataSourcesRuntimeControl{sync.Mutex{}, cancel, shutdown}
}

func (r *RuntimeManager) UpdateDataSource(ctx context.Context, ds datasources.DataSource) {
	r.StopDataSource(ds.Id())
	r.RemoveDataSource(ds.Id())
	r.AddDataSource(ds)
	r.Run(ctx, ds.Id())
}

func (r *RuntimeManager) RestartDataSource(ctx context.Context, id shared.CommonId) {
	r.Logger.Infof("Restarting datasource %s", id.String())
	r.StopDataSource(id)
	r.Run(ctx, id)
}

func (r *RuntimeManager) StopDataSource(id shared.CommonId) {
	start := time.Now()

	r.Logger.Infof("Shutdown datapoint runtime %s", id.String())

	shutdown := r.dataSourcesRuntimeControl[id]

	if shutdown == nil {
		return
	}

	shutdown.Lock()
	defer shutdown.Unlock()

	shutdown.shutdown()
	<-shutdown.confirmShutdown

	delete(r.dataSourcesRuntimeControl, id)

	r.Logger.Infof("Shutdown completed of data source id: %s take %d Milliseconds", id.String(), time.Since(start).Milliseconds())
}

func (r *RuntimeManager) StopAll() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(len(r.dataSourcesRuntimeControl))

	for id := range r.dataSourcesRuntimeControl {
		id := id
		go func() {
			r.StopDataSource(id)
			wg.Done()
		}()
	}

	wg.Wait()
}
