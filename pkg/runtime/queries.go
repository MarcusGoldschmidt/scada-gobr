package runtime

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"sync"
)

type DataSourceRuntimeManagerStatus struct {
	Id     shared.CommonId `json:"id"`
	Name   string          `json:"name"`
	Status string          `json:"status"`
	Error  string          `json:"error"`
}

func newDataSourceRuntimeManagerStatus(manager datasources.DataSourceRuntimeManager) *DataSourceRuntimeManagerStatus {
	err := manager.GetError()
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	return &DataSourceRuntimeManagerStatus{
		Id:     manager.Id(),
		Name:   manager.Name(),
		Status: manager.Status().String(),
		Error:  errorMsg,
	}
}

func (r *Manager) GetAllDataSources() []*DataSourceRuntimeManagerStatus {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	response := make([]*DataSourceRuntimeManagerStatus, len(r.dataSources))

	wg := sync.WaitGroup{}
	wg.Add(len(r.dataSources))
	i := 0
	for _, manager := range r.dataSources {
		go func(index int, controller *dataSourceController) {
			controller.mutex.RLock()
			defer controller.mutex.RUnlock()
			defer wg.Done()

			response[index] = newDataSourceRuntimeManagerStatus(controller.dataSourceRuntimeManager)
		}(i, manager)

		i++
	}

	wg.Wait()

	return response
}
