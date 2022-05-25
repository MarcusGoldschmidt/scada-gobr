package datasources

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/shared"
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

type HttpRequestDataSource struct {
	*DataSourceCommon
	period     time.Duration
	dataPoints []*HttpRequestDataPoint

	baseUrl      string
	encoding     string
	method       string
	headers      map[string]string
	bodyTemplate *string

	forEachDataPoint bool
}

type HttpRequestDataSourceRuntime struct {
	id          shared.CommonId
	dataSource  HttpRequestDataSource
	persistence persistence.DataPointPersistence
	client      *http.Client
}

func (r *HttpRequestDataSource) AddDataPoint(dp *HttpRequestDataPoint) {
	r.dataPoints = append(r.dataPoints, dp)
}

func (r HttpRequestDataSource) Id() shared.CommonId {
	return r.id
}

func (r HttpRequestDataSource) Name() string {
	return r.name
}

func (r HttpRequestDataSource) IsEnable() bool {
	return r.isEnable
}

func (r HttpRequestDataSource) GetDataPoints() []Datapoint {
	datapoint := make([]Datapoint, len(r.dataPoints))
	for i, v := range r.dataPoints {
		datapoint[i] = Datapoint(v)
	}
	return datapoint
}

func (r HttpRequestDataSource) CreateRuntime(ctx context.Context, p persistence.DataPointPersistence) (DataSourceRuntime, error) {
	rt := HttpRequestDataSourceRuntime{uuid.New(), r, p, &http.Client{}}
	return rt, nil
}

func (c HttpRequestDataSourceRuntime) GetDataSource() DataSource {
	return c.dataSource
}

func (c HttpRequestDataSourceRuntime) Run(ctx context.Context, shutdownCompleteChan chan shared.CommonId) error {
	defer func() {
		shutdownCompleteChan <- c.id
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(c.dataSource.period):
			series, err := c.QueryDatabase(ctx)
			if err != nil {
				return err
			}
			err = c.persistence.AddDataPointValues(ctx, series)
			if err != nil {
				return err
			}
		}
	}
}

func (c *HttpRequestDataSourceRuntime) QueryDatabase(ctx context.Context) ([]*shared.IdSeries, error) {
	dict := make(map[string]*HttpRequestDataPoint)

	for _, point := range c.dataSource.dataPoints {
		dict[point.rowIdentifier] = point
	}

	bodyTemplate := ""

	if c.dataSource.forEachDataPoint && c.dataSource.bodyTemplate != nil {
		resp, err := c.parseBodyTemplate()
		if err != nil {
			return nil, err
		}
		bodyTemplate = *resp
	}

	req, err := http.NewRequest(c.dataSource.method, c.dataSource.baseUrl, bytes.NewBufferString(bodyTemplate))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData []map[string]string

	if c.dataSource.encoding == "XML" {
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

func (c HttpRequestDataSourceRuntime) parseBodyTemplate() (*string, error) {

	if c.dataSource.bodyTemplate == nil {
		return nil, errors.New("do not found body template")
	}

	tmpl, err := template.New("test").Parse(*c.dataSource.bodyTemplate)
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
