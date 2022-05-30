package datasources

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type HttpRequestDataPoint struct {
	*DataPointCommon
	rowIdentifier string
	dateFormat    string
}

func (r HttpRequestDataPoint) Id() shared.CommonId {
	return r.id
}

func (r HttpRequestDataPoint) Name() string {
	return r.name
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

func NewHttpRequestWorker(dataSourceId shared.CommonId) *HttpRequestWorker {
	return &HttpRequestWorker{dataSourceId: dataSourceId}
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

func (c *HttpRequestWorker) QueryDatabase(ctx context.Context) ([]*shared.IdSeries, error) {
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

	result := make([]*shared.IdSeries, 0)

	for _, data := range responseData {
		value, err := strconv.ParseFloat(data["value"], 64)

		if err != nil {
			return nil, err
		}

		if point, ok := dict["name"]; ok {
			timestamp, err := time.Parse(point.dateFormat, data["timestamp"])
			if err != nil {
				return nil, err
			}
			result = append(result, shared.NewIdSeries(point.id, value, timestamp))
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
