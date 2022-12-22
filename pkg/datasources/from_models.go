package datasources

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
)

func ModelDataPointToCommon[T any](dp *models.DataPoint, cb func(*DataPointCommon, *models.DataPoint) (*T, error)) (*T, error) {
	common := NewDataPointCommon(dp.Id, dp.Name, dp.IsEnable)

	return cb(common, dp)
}

func ModelDataPointToSql(dp *models.DataPoint) (*SqlDataPoint, error) {
	fun := func(common *DataPointCommon, data *models.DataPoint) (*SqlDataPoint, error) {
		dsTypeData, err := shared.FromJson[models.DataPointTypeSql](dp.Data)

		if err != nil {
			return nil, err
		}

		return NewSqlDataPoint(common, dsTypeData.RowIdentifier), nil
	}

	return ModelDataPointToCommon[SqlDataPoint](dp, fun)
}

func ModelDataPointToRandom(dp *models.DataPoint) (*RandomValueDataPoint, error) {
	fun := func(common *DataPointCommon, data *models.DataPoint) (*RandomValueDataPoint, error) {
		dsTypeData, err := shared.FromJson[models.DataPointTypeRandomValue](dp.Data)

		if err != nil {
			return nil, err
		}

		return NewRandomValueDataPoint(common, dsTypeData.InitialValue, dsTypeData.EndValue), err
	}

	return ModelDataPointToCommon[RandomValueDataPoint](dp, fun)
}

func ModelDataPointToHttpRequest(dp *models.DataPoint) (*HttpRequestDataPoint, error) {
	fun := func(common *DataPointCommon, data *models.DataPoint) (*HttpRequestDataPoint, error) {
		dsTypeData, err := shared.FromJson[models.DataPointTypeHttpRequest](dp.Data)

		if err != nil {
			return nil, err
		}

		return NewHttpRequestDataPoint(common, dsTypeData.RowIdentifier, dsTypeData.DateFormat), err
	}

	return ModelDataPointToCommon[HttpRequestDataPoint](dp, fun)
}

func ModelDataPointToHttpServer(dp *models.DataPoint) (*HttpServerDataPoint, error) {
	fun := func(common *DataPointCommon, data *models.DataPoint) (*HttpServerDataPoint, error) {
		dsTypeData, err := shared.FromJson[models.DataPointTypeHttpServer](dp.Data)

		if err != nil {
			return nil, err
		}

		return NewHttpServerDataPoint(common, dsTypeData.RowIdentifier, dsTypeData.DateFormat), err
	}

	return ModelDataPointToCommon[HttpServerDataPoint](dp, fun)
}
