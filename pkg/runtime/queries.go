package runtime

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
)

type DataSourceRuntimeManagerStatus struct {
	Id     shared.CommonId              `json:"id"`
	Name   string                       `json:"name"`
	Status datasources.DataSourceStatus `json:"status"`
	Error  string                       `json:"error"`
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
		Status: manager.Status(),
		Error:  errorMsg,
	}
}

func (r *RuntimeManager) GetAllDataSources() []*DataSourceRuntimeManagerStatus {

	response := make([]*DataSourceRuntimeManagerStatus, len(r.dataSources))

	for _, manager := range r.dataSources {
		response = append(response, newDataSourceRuntimeManagerStatus(manager))
	}

	return response
}
