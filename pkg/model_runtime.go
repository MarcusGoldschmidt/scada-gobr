package pkg

import (
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"net/http"
)

func DataPointToRuntimeFunc[T any](dp *models.DataPoint, cb func(*datasources.DataPointCommon, *models.DataPoint) (*T, error)) (*T, error) {
	common := datasources.NewDataPointCommon(dp.Id, dp.Name, dp.IsEnable)

	return cb(common, dp)
}

func DataPointToRuntimeSql(dp *models.DataPoint) (*datasources.SqlDataPoint, error) {
	fun := func(common *datasources.DataPointCommon, data *models.DataPoint) (*datasources.SqlDataPoint, error) {
		dsTypeData, err := shared.FromJson[models.DataPointTypeSql](dp.Data)

		if err != nil {
			return nil, err
		}

		return datasources.NewSqlDataPoint(common, dsTypeData.RowIdentifier), nil
	}

	return DataPointToRuntimeFunc[datasources.SqlDataPoint](dp, fun)
}

func DataPointToRuntimeRandom(dp *models.DataPoint) (*datasources.RandomValueDataPoint, error) {

	fun := func(common *datasources.DataPointCommon, data *models.DataPoint) (*datasources.RandomValueDataPoint, error) {
		dsTypeData, err := shared.FromJson[models.DataPointTypeRandomValue](dp.Data)

		if err != nil {
			return nil, err
		}

		return datasources.NewRandomValueDataPoint(common, dsTypeData.InitialValue, dsTypeData.EndValue), err
	}

	return DataPointToRuntimeFunc[datasources.RandomValueDataPoint](dp, fun)
}

func DataPointToRuntimeHttpRequest(dp *models.DataPoint) (*datasources.HttpRequestDataPoint, error) {

	fun := func(common *datasources.DataPointCommon, data *models.DataPoint) (*datasources.HttpRequestDataPoint, error) {
		dsTypeData, err := shared.FromJson[models.DataPointTypeHttpRequest](dp.Data)

		if err != nil {
			return nil, err
		}

		return datasources.NewHttpRequestDataPoint(common, dsTypeData.RowIdentifier, dsTypeData.DateFormat), err
	}

	return DataPointToRuntimeFunc[datasources.HttpRequestDataPoint](dp, fun)
}

func DataPointToRuntimeHttpServer(dp *models.DataPoint) (*datasources.HttpServerDataPoint, error) {
	fun := func(common *datasources.DataPointCommon, data *models.DataPoint) (*datasources.HttpServerDataPoint, error) {
		dsTypeData, err := shared.FromJson[models.DataPointTypeHttpServer](dp.Data)

		if err != nil {
			return nil, err
		}

		return datasources.NewHttpServerDataPoint(common, dsTypeData.RowIdentifier, dsTypeData.DateFormat), err
	}

	return DataPointToRuntimeFunc[datasources.HttpServerDataPoint](dp, fun)
}

func DataSourceToRuntimeManager(scada *Scadagobr, ds *models.DataSource) (datasources.DataSourceRuntimeManager, error) {
	logger := scada.RuntimeManager.CreateLogger(ds.Id, ds.Name)
	runtimeManager := datasources.NewDataSourceRuntimeManagerCommon(ds.Id, ds.Name, logger)

	ds.FilterAvailableDataPoints()

	var worker datasources.DataSourceWorker

	// TODO: refactor this
	// Should I use glob?
	switch ds.Type {
	case models.Sql:
		dsTypeData, err := shared.FromJson[models.DataSourceTypeSql](ds.Data)

		if err != nil {
			return nil, err
		}

		dataPoints := make([]*datasources.SqlDataPoint, 0)
		for _, point := range ds.DataPoints {
			sql, err := DataPointToRuntimeSql(point)
			if err != nil {
				return nil, err
			}
			dataPoints = append(dataPoints, sql)
		}

		worker = datasources.NewSqlWorker(
			dsTypeData.Period,
			dataPoints,
			dsTypeData.Driver,
			dsTypeData.Query,
			dsTypeData.ConnectionString,
			ds.Id,
			scada.dataPointPersistence,
		)
		break
	case models.RandomValue:
		dsTypeData, err := shared.FromJson[models.DataSourceTypeRandomValue](ds.Data)
		if err != nil {
			return nil, err
		}

		dataPoints := make([]*datasources.RandomValueDataPoint, 0)
		for _, point := range ds.DataPoints {
			random, err := DataPointToRuntimeRandom(point)
			if err != nil {
				return nil, err
			}
			dataPoints = append(dataPoints, random)
		}

		worker = datasources.NewRandomValueWorker(ds.Id, dsTypeData.Period, dataPoints, scada.dataPointPersistence)
		break
	case models.HttpRequest:
		dsTypeData, err := shared.FromJson[models.DataSourceTypeHttpRequest](ds.Data)
		if err != nil {
			return nil, err
		}

		dataPoints := make([]*datasources.HttpRequestDataPoint, 0)
		for _, point := range ds.DataPoints {
			httpRequestDp, err := DataPointToRuntimeHttpRequest(point)
			if err != nil {
				return nil, err
			}
			dataPoints = append(dataPoints, httpRequestDp)
		}

		response := datasources.HttpRequestWorker{
			Period:           dsTypeData.Period,
			DataPoints:       dataPoints,
			Method:           dsTypeData.Method,
			Persistence:      scada.dataPointPersistence,
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
		dsTypeData, err := shared.FromJson[models.DataSourceTypeHttpServer](ds.Data)
		if err != nil {
			return nil, err
		}

		dataPoints := make([]*datasources.HttpServerDataPoint, len(ds.DataPoints))
		for i, point := range ds.DataPoints {
			httpRequestDp, err := DataPointToRuntimeHttpServer(point)
			if err != nil {
				return nil, err
			}
			dataPoints[i] = httpRequestDp
		}

		response := datasources.HttpServerWorker{
			DataPoints:   dataPoints,
			Persistence:  scada.dataPointPersistence,
			User:         dsTypeData.User,
			PasswordHash: dsTypeData.PasswordHash,
			Router:       scada.internalRoute,
			AtmDone:      0,
			Endpoint:     dsTypeData.Endpoint,
		}
		response.SetDataSourceId(ds.Id)

		worker = &response
	}

	if worker == nil {
		return nil, errors.New("worker not found")
	}

	runtimeManager.WithWorker(worker)

	return runtimeManager, nil
}
