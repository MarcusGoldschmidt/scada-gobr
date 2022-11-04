package pkg

import (
	"encoding/json"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
)

func GetDataSourcesRuntime(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	s.respondJsonOk(r.Context(), w, s.RuntimeManager.GetAllDataSources())
}

func GetDataSourceTypesHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	s.respondJsonOk(r.Context(), w, datasources.DataSourceTypes)
}

func GetDataSourcesHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sources, err := s.dataSourcePersistence.GetDataSources(ctx)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	for _, source := range sources {
		err := json.Unmarshal(source.Data, &source.TypeData)
		if err != nil {
			s.respondError(ctx, w, err)
			return
		}
	}

	s.respondJsonOk(ctx, w, sources)
}

func GetDataSourceByIdHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	datasource, err := s.dataSourcePersistence.GetDataSourceById(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = json.Unmarshal(datasource.Data, &datasource.TypeData)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, datasource)
}

type createDataSource struct {
	Name string                `json:"name"`
	Data map[string]any        `json:"data"`
	Type models.DataSourceType `json:"type"`
}

func CreateDataSourceHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	command, err := server.ReadJson[createDataSource](r)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	existDataSource, err := s.dataSourcePersistence.GetDataSourceByName(ctx, command.Name)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if existDataSource != nil {
		s.respondBadRequest(ctx, w, errors.New("datasource already exist"))
		return
	}

	datasource := &models.DataSource{
		Id:   uuid.New(),
		Name: command.Name,
		Type: command.Type,
	}

	err = parseDataSourceData(datasource, command.Data)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.dataSourcePersistence.CreateDataSource(ctx, datasource)
	if err != nil {
		return
	}

	s.respondJsonOk(ctx, w, command)
}

func EditDataSourceHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	claims, err := auth.GetUserFromContext(ctx)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if !claims.Admin {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	command, err := server.ReadJson[createDataSource](r)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	datasource, err := s.dataSourcePersistence.GetDataSourceById(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	datasource.Name = command.Name

	err = parseDataSourceData(datasource, command.Data)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, command)
}

func DeleteDataSourceHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	datapoints, err := s.dataSourcePersistence.GetDataPoints(ctx, id)

	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.dataSourcePersistence.DeleteDataSource(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(datapoints) + 1)

	go func() {
		defer wg.Done()
		_ = s.RuntimeManager.StopDataSource(ctx, id)
	}()

	for _, datapoint := range datapoints {
		go func(datapoint *models.DataPoint) {
			defer wg.Done()
			_ = s.dataPointPersistence.DeleteDataPointValueById(ctx, datapoint.Id)
		}(datapoint)
	}

	wg.Wait()

	w.WriteHeader(http.StatusOK)
}

func parseDataSourceData(dataSource *models.DataSource, data map[string]any) error {
	var marshal []byte
	var err error

	switch dataSource.Type {
	case models.Sql:
		marshal, err = shared.ValidateDataSourceType[models.DataSourceTypeSql](data)
	case models.HttpRequest:
		marshal, err = shared.ValidateDataSourceType[models.DataSourceTypeHttpRequest](data)
	case models.HttpServer:
		marshal, err = shared.ValidateDataSourceType[models.DataSourceTypeHttpServer](data)
	case models.RandomValue:
		if value, ok := data["period"].(string); ok {
			duration, err := time.ParseDuration(value)

			if err != nil {
				return err
			}

			model := models.DataSourceTypeRandomValue{
				Period: duration,
			}

			marshal, err = json.Marshal(model)
		}
	default:
		err = errors.New("unknown datasource type")
	}

	dataSource.Data = marshal

	return err
}
