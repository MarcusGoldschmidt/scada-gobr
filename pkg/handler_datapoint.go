package pkg

import (
	"encoding/json"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func GetDataPointByIdHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	points, err := s.dataPointPersistence.GetDataPointById(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, points)
}

func GetDataPointsHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	points, err := s.dataSourcePersistence.GetDataPoints(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, points)
}

type createDataPoint struct {
	Name       string         `json:"name"`
	IsEnable   bool           `json:"isEnable"`
	Unit       string         `json:"unit"`
	PurgeAfter *time.Duration `json:"purgeAfter"`
	Data       map[string]any `json:"data"`
}

func CreateDataPointHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	command, err := server.ValidateFromBody[createDataPoint](r)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	ds, err := s.dataSourcePersistence.GetDataSourceById(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	dataPoint := &models.DataPoint{
		Id:           uuid.New(),
		Name:         command.Name,
		DataSourceId: id,
		IsEnable:     command.IsEnable,
		PurgeAfter:   command.PurgeAfter,
		Unit:         command.Unit,
		Type:         ds.Type,
	}

	err = parseDataPointData(dataPoint, command.Data)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.dataPointPersistence.CreateDataPoint(ctx, dataPoint)
	if err != nil {
		return
	}

	err = s.UpdateDataSource(ctx, ds.Id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = json.Unmarshal(dataPoint.Data, &dataPoint.TypeData)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, dataPoint)
}

func EditDataPointHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	dataPointId, err := uuid.Parse(vars["dataPointId"])
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

	command, err := server.ValidateFromBody[createDataPoint](r)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	dataPoint, err := s.dataPointPersistence.GetDataPointById(ctx, dataPointId)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if dataPoint.DataSourceId != id {
		s.respondError(ctx, w, errors.New("invalid data point id"))
		return
	}

	dataPoint.Name = command.Name
	dataPoint.IsEnable = command.IsEnable
	dataPoint.Unit = command.Unit
	dataPoint.PurgeAfter = command.PurgeAfter

	err = parseDataPointData(dataPoint, command.Data)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.dataPointPersistence.CreateDataPoint(ctx, dataPoint)
	if err != nil {
		return
	}

	err = json.Unmarshal(dataPoint.Data, &dataPoint.TypeData)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.dataPointPersistence.UpdateDataPoint(ctx, dataPoint)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.UpdateDataSource(ctx, dataPoint.DataSourceId)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, dataPoint)
}

func DeleteDataPointHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	dataPointId, err := uuid.Parse(vars["dataPointId"])
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

	err = s.dataPointPersistence.DeleteDataPoint(ctx, id, dataPointId)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.UpdateDataSource(ctx, id)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	err = s.dataPointPersistence.DeleteDataPointValueById(ctx, dataPointId)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func parseDataPointData(dataPoint *models.DataPoint, data map[string]any) error {
	var marshal []byte
	var err error

	switch dataPoint.Type {
	case models.Sql:
		marshal, err = shared.ValidateDataSourceType[models.DataPointTypeSql](data)
	case models.HttpRequest:
		marshal, err = shared.ValidateDataSourceType[models.DataPointTypeHttpRequest](data)
	case models.HttpServer:
		marshal, err = shared.ValidateDataSourceType[models.DataPointTypeHttpServer](data)
	case models.RandomValue:
		marshal, err = shared.ValidateDataSourceType[models.DataPointTypeRandomValue](data)
	default:
		err = errors.New("unknown datasource type")
	}

	dataPoint.Data = marshal

	return err
}
