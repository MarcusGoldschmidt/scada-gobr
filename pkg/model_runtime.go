package pkg

import (
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"net/http"
)

func parseDataPoints[T any](dataPoints []*models.DataPoint, transform func(dp *models.DataPoint) (*T, error)) ([]*T, error) {
	result := make([]*T, 0)

	for _, point := range dataPoints {
		newPoint, err := transform(point)
		if err != nil {
			return nil, err
		}
		result = append(result, newPoint)
	}

	return result, nil
}

func DataSourceToRuntimeManager(scada *Scadagobr, ds *models.DataSource) (datasources.DataSourceRuntimeManager, error) {
	logger := scada.RuntimeManager.CreateLogger(ds.Id, ds.Name)

	availableDataPoints := ds.FilterAvailableDataPoints()

	var worker datasources.DataSourceWorker

	// TODO: refactor this
	// Should I use gob?
	switch ds.Type {
	case models.Sql:
		dsTypeData, err := models.ParseDataSourceTypeData[models.DataSourceTypeSql](ds)

		if err != nil {
			return nil, err
		}

		dataPoints, err := parseDataPoints[datasources.SqlDataPoint](availableDataPoints, datasources.ModelDataPointToSql)
		if err != nil {
			return nil, err
		}

		worker = datasources.NewSqlWorker(
			dsTypeData.Period,
			dataPoints,
			dsTypeData.Driver,
			dsTypeData.Query,
			dsTypeData.ConnectionString,
			ds.Id,
			scada.DataPointPersistence,
		)
		break
	case models.RandomValue:
		dsTypeData, err := models.ParseDataSourceTypeData[models.DataSourceTypeRandomValue](ds)
		if err != nil {
			return nil, err
		}

		dataPoints, err := parseDataPoints[datasources.RandomValueDataPoint](availableDataPoints, datasources.ModelDataPointToRandom)
		if err != nil {
			return nil, err
		}

		worker = datasources.NewRandomValueWorker(ds.Id, dsTypeData.Period, dataPoints, scada.DataPointPersistence)
		break
	case models.HttpRequest:
		dsTypeData, err := models.ParseDataSourceTypeData[models.DataSourceTypeHttpRequest](ds)
		if err != nil {
			return nil, err
		}

		dataPoints, err := parseDataPoints[datasources.HttpRequestDataPoint](availableDataPoints, datasources.ModelDataPointToHttpRequest)
		if err != nil {
			return nil, err
		}

		response := datasources.HttpRequestWorker{
			Period:           dsTypeData.Period,
			DataPoints:       dataPoints,
			Method:           dsTypeData.Method,
			Persistence:      scada.DataPointPersistence,
			BaseUrl:          dsTypeData.BaseUrl,
			Client:           &http.Client{},
			BodyTemplate:     dsTypeData.BodyTemplate,
			Encoding:         dsTypeData.Encoding,
			ForEachDataPoint: dsTypeData.ForEachDataPoint,
			Headers:          dsTypeData.Headers,
		}
		response.SetDataSourceId(ds.Id)

		worker = &response
		break
	case models.HttpServer:
		dsTypeData, err := models.ParseDataSourceTypeData[models.DataSourceTypeHttpServer](ds)
		if err != nil {
			return nil, err
		}

		dataPoints, err := parseDataPoints[datasources.HttpServerDataPoint](availableDataPoints, datasources.ModelDataPointToHttpServer)
		if err != nil {
			return nil, err
		}

		response := datasources.HttpServerWorker{
			DataPoints:   dataPoints,
			Persistence:  scada.DataPointPersistence,
			User:         dsTypeData.User,
			PasswordHash: dsTypeData.PasswordHash,
			Router:       scada.InternalRouter,
			AtmDone:      0,
			Endpoint:     dsTypeData.Endpoint,
			DataSourceId: ds.Id,
		}

		worker = &response
	default:
		return nil, errors.New("unknown data source type: " + string(ds.Type))
	}

	if worker == nil {
		return nil, errors.New("worker not found")
	}

	runtimeManager := datasources.NewDataSourceRuntimeManagerCommon(ds.Id, ds.Name, worker, logger)

	return runtimeManager, nil
}
