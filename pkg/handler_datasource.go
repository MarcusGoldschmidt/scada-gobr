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
)

func GetDataSourcesRuntime(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	s.respondJsonOk(w, s.RuntimeManager.GetAllDataSources())
}

func GetDataSourceTypesHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	s.respondJsonOk(w, datasources.DataSourceTypes)
}

func GetDataSourcesHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sources, err := s.dataSourcePersistence.GetDadaSources(ctx)
	if err != nil {
		s.respondError(w, err)
		return
	}

	for _, source := range sources {
		err := json.Unmarshal(source.Data, &source.TypeData)
		if err != nil {
			s.respondError(w, err)
			return
		}
	}

	s.respondJsonOk(w, sources)
}

type createDataSource struct {
	Name string         `json:"name"`
	Data map[string]any `json:"data"`
}

func CreateDataSourceHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dataSourceId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(w, err)
		return
	}

	command, err := server.ReadJson[createDataSource](r)
	if err != nil {
		s.respondError(w, err)
		return
	}

	ds, err := s.dataSourcePersistence.GetDadaSourceById(ctx, dataSourceId)
	if err != nil {
		s.respondError(w, err)
		return
	}

	datasource := &models.DataSource{
		Id:   uuid.New(),
		Name: command.Name,
		Type: ds.Type,
	}

	err = parseDataSourceData(datasource, command.Data)
	if err != nil {
		s.respondError(w, err)
		return
	}

	err = s.dataSourcePersistence.CreateDataSource(ctx, datasource)
	if err != nil {
		return
	}

	s.respondJsonOk(w, command)
}

func EditDataSourceHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(w, err)
		return
	}

	claims, err := auth.GetUserFromContext(ctx)
	if err != nil {
		s.respondError(w, err)
		return
	}

	if !claims.Admin {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	command, err := server.ReadJson[createDataSource](r)
	if err != nil {
		s.respondError(w, err)
		return
	}

	datasource, err := s.dataSourcePersistence.GetDadaSourceById(ctx, id)
	if err != nil {
		s.respondError(w, err)
		return
	}

	datasource.Name = command.Name

	err = parseDataSourceData(datasource, command.Data)
	if err != nil {
		s.respondError(w, err)
		return
	}

	s.respondJsonOk(w, command)
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
		marshal, err = shared.ValidateDataSourceType[models.DataSourceTypeRandomValue](data)
	default:
		err = errors.New("unknown datasource type")
	}

	dataSource.Data = marshal

	return err
}
