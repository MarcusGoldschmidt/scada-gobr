package datasources

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type HttpRequestDataPoint struct {
	common        *DataPointCommon
	rowIdentifier string
	dateFormat    string
}

func NewHttpRequestDataPoint(common *DataPointCommon, rowIdentifier string, dateFormat string) *HttpRequestDataPoint {
	return &HttpRequestDataPoint{common: common, rowIdentifier: rowIdentifier, dateFormat: dateFormat}
}

func (r HttpRequestDataPoint) Id() shared.CommonId {
	return r.common.Id
}

func (r HttpRequestDataPoint) Name() string {
	return r.common.Name
}

type HttpRequestWorker struct {
	Period     time.Duration
	DataPoints []*HttpRequestDataPoint

	BaseUrl      string
	Encoding     string
	Method       string
	Headers      map[string]string
	BodyTemplate *string

	ForEachDataPoint bool

	Persistence persistence.DataPointPersistence
	Client      *http.Client

	dataSourceId shared.CommonId
}

func (c *HttpRequestWorker) SetDataSourceId(dataSourceId shared.CommonId) {
	c.dataSourceId = dataSourceId
}

func (c *HttpRequestWorker) DataSourceId() shared.CommonId {
	return c.dataSourceId
}

func (c *HttpRequestWorker) Run(ctx context.Context, errorChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.Period):
			series, err := c.QueryDatabase(ctx)
			if err != nil {
				errorChan <- err
				return
			}
			err = c.Persistence.AddDataPointValues(ctx, series)
			if err != nil {
				errorChan <- err
				return
			}
		}
	}
}

func (c *HttpRequestWorker) QueryDatabase(ctx context.Context) ([]*models.DataSeries, error) {
	dict := make(map[string]*HttpRequestDataPoint)

	for _, point := range c.DataPoints {
		dict[point.rowIdentifier] = point
	}

	bodyTemplate := ""

	if c.ForEachDataPoint && c.BodyTemplate != nil {
		resp, err := c.parseBodyTemplate()
		if err != nil {
			return nil, err
		}
		bodyTemplate = *resp
	}

	req, err := http.NewRequest(c.Method, c.BaseUrl, bytes.NewBufferString(bodyTemplate))
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData []map[string]string

	if c.Encoding == "XML" {
		err := xml.Unmarshal(responseBody, &responseData)
		if err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal(responseBody, &responseData)
		if err != nil {
			return nil, err
		}
	}

	result := make([]*models.DataSeries, 0)

	for _, data := range responseData {
		value, err := strconv.ParseFloat(data["value"], 64)

		if err != nil {
			return nil, err
		}

		if point, ok := dict["Name"]; ok {
			timestamp, err := time.Parse(point.dateFormat, data["timestamp"])
			if err != nil {
				return nil, err
			}
			result = append(result, models.NewDataSeries(timestamp, value, point.Id()))
		}
	}

	return result, nil
}

func (c *HttpRequestWorker) parseBodyTemplate() (*string, error) {

	if c.BodyTemplate == nil {
		return nil, errors.New("do not found body template")
	}

	tmpl, err := template.New("test").Parse(*c.BodyTemplate)
	if err != nil {
		return nil, err
	}

	dict := map[string]string{}

	buff := bytes.Buffer{}

	err = tmpl.Execute(&buff, dict)
	if err != nil {
		return nil, err
	}

	response := buff.String()

	return &response, nil
}
