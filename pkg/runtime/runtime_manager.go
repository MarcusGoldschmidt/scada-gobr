package runtime

import (
	"context"
	"fmt"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/buffers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events/topics"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/logger"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"io"
	"os"
	"sync"
)

type ManagerOptions struct {
	MaxRuntimeRetry int
}

type dataSourceController struct {
	dataSourceRuntimeManager datasources.DataSourceRuntimeManager
	mutex                    sync.RWMutex
	logs                     *buffers.MaxBuffer
}

func newDataSourceController(dataSourceRuntimeManager datasources.DataSourceRuntimeManager) *dataSourceController {
	return &dataSourceController{
		dataSourceRuntimeManager,
		sync.RWMutex{},
		nil,
	}
}

type Manager struct {
	Logger logger.Logger
	mutex  sync.RWMutex

	dataSources map[shared.CommonId]*dataSourceController

	options ManagerOptions

	persistence  persistence.DataPointPersistence
	timeProvider providers.TimeProvider

	eventsManager events.HubManager
}

func NewRuntimeManager(logger logger.Logger, persistence persistence.DataPointPersistence, eventsManager events.HubManager, provider providers.TimeProvider) *Manager {
	return &Manager{
		Logger:        logger,
		mutex:         sync.RWMutex{},
		dataSources:   make(map[shared.CommonId]*dataSourceController),
		persistence:   persistence,
		timeProvider:  provider,
		options:       ManagerOptions{MaxRuntimeRetry: 5},
		eventsManager: eventsManager,
	}
}

func (r *Manager) WithOptions(opt ManagerOptions) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.options = opt
}

func (r *Manager) AddDataSourceManager(sources ...datasources.DataSourceRuntimeManager) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, source := range sources {
		r.dataSources[source.Id()] = newDataSourceController(source)
	}
}

func (r *Manager) RemoveDataSource(ctx context.Context, id shared.CommonId) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.dataSources, id)
	r.NotifyUpdated(ctx)
}

func (r *Manager) Run(ctx context.Context, id shared.CommonId) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	dataSource, ok := r.dataSources[id]

	if !ok {
		return fmt.Errorf("data source %s not found", id.String())
	}

	dataSource.mutex.Lock()
	defer dataSource.mutex.Unlock()

	dataSource.dataSourceRuntimeManager.WithNotificationEachStatus(func(ctx context.Context, manager datasources.DataSourceRuntimeManager) {
		r.NotifyUpdated(ctx)
	})

	err := dataSource.dataSourceRuntimeManager.Run(ctx)
	if err != nil {
		r.Logger.Errorf("datasource runtime %s stopped with error: %s", dataSource.dataSourceRuntimeManager.Name(), err.Error())
		return err
	}

	r.NotifyUpdated(ctx)

	return nil
}

func (r *Manager) RunAll(ctx context.Context) error {
	for id, _ := range r.dataSources {
		err := r.Run(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Manager) UpdateDataSource(ctx context.Context, ds datasources.DataSourceRuntimeManager) error {
	err := r.StopDataSource(ctx, ds.Id())

	if err != nil {
		return err
	}

	r.RemoveDataSource(ctx, ds.Id())
	r.AddDataSourceManager(ds)
	// TODO: parse trace id
	return r.Run(context.Background(), ds.Id())
}

func (r *Manager) RestartDataSource(ctx context.Context, id shared.CommonId) error {
	r.Logger.Infof("Restarting datasource %s", id.String())
	err := r.StopDataSource(ctx, id)
	if err != nil {
		return err
	}
	return r.Run(ctx, id)
}

func (r *Manager) StopDataSource(ctx context.Context, id shared.CommonId) error {
	r.Logger.Infof("Shutdown datapoint runtime %s", id.String())

	datasourceManager, ok := r.dataSources[id]

	if !ok {
		r.Logger.Warningf("Shutdown datapoint runtime %s not found", id.String())
		return nil
	}

	datasourceManager.mutex.Lock()
	defer datasourceManager.mutex.Unlock()

	err := datasourceManager.dataSourceRuntimeManager.Stop(ctx)
	if err != nil {
		return err
	}

	r.NotifyUpdated(ctx)

	return nil
}

func (r *Manager) StopAll(ctx context.Context) {
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

func (r *Manager) CreateLogger(id shared.CommonId, name string) logger.Logger {
	dataSource, ok := r.dataSources[id]

	if !ok {
		return logger.NewSimpleLogger(id.String()+"-"+name, os.Stdout)
	}
	var bufferSize *buffers.MaxBuffer
	if dataSource.logs != nil {
		bufferSize = dataSource.logs
	}

	bufferSize = buffers.NewMaxBuffer(buffers.MB)

	dataSource.logs = bufferSize

	logOutput := io.MultiWriter(os.Stdout, bufferSize)

	return logger.NewSimpleLogger(id.String()+"-"+name, logOutput)
}

func (r *Manager) NotifyUpdated(ctx context.Context) {
	r.eventsManager.SendMessage(ctx, topics.RuntimeManagerUpdated, nil)
}
