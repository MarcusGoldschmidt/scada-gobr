package datasources

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type HttpServerDataPoint struct {
	common        *DataPointCommon
	rowIdentifier string
	dateFormat    string
}

func NewHttpServerDataPoint(common *DataPointCommon, rowIdentifier string, dateFormat string) *HttpServerDataPoint {
	return &HttpServerDataPoint{common: common, rowIdentifier: rowIdentifier, dateFormat: dateFormat}
}

func (r HttpServerDataPoint) Id() shared.CommonId {
	return r.common.Id
}

func (r HttpServerDataPoint) Name() string {
	return r.common.Name
}

type HttpServerWorker struct {
	DataPoints []*HttpServerDataPoint

	Router       *server.Router
	Endpoint     string
	User         string
	PasswordHash string

	Persistence persistence.DataPointPersistence
	AtmDone     int32

	dataSourceId shared.CommonId
}

func (c *HttpServerWorker) SetDataSourceId(dataSourceId shared.CommonId) {
	c.dataSourceId = dataSourceId
}

func (c *HttpServerWorker) WithDataSourceId(dataSourceId shared.CommonId) {
	c.dataSourceId = dataSourceId
}

func (c *HttpServerWorker) DataSourceId() shared.CommonId {
	return c.dataSourceId
}

type request struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func (c *HttpServerWorker) Run(ctx context.Context, errorChan chan<- error) {
	// TODO: verify best number
	channel := make(chan []*request, 1024)

	c.GerOrAddRoute(channel)
	defer c.RemoveRoute()
	defer close(channel)
	defer atomic.AddInt32(&c.AtmDone, 1)

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-channel:
			dict := make(map[string]shared.CommonId)

			for _, point := range c.DataPoints {
				dict[point.rowIdentifier] = point.Id()
			}

			series := make([]*models.DataSeries, 0)
			for _, d := range data {
				if id, ok := dict[d.Name]; ok {
					series = append(series, models.NewDataSeries(d.Timestamp, d.Value, id))
				}
			}

			err := c.Persistence.AddDataPointValues(ctx, series)
			if err != nil {
				errorChan <- err
				return
			}
		}
	}
}

func (c *HttpServerWorker) GerOrAddRoute(channel chan []*request) http.Handler {
	if data, ok := c.Router.VerifyMatch(c.Endpoint); ok {
		return data
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		basicAuth := r.Header["Authorization"][0][6:]

		split := strings.Split(basicAuth, ":")

		hash := sha256.Sum256([]byte(split[1]))

		if split[0] != c.User && c.PasswordHash != string(hash[:]) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if c.AtmDone == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("datasource runtime is done, try again later"))
			return
		}

		var data []*request

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			setupError(w, err)
			return
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			setupError(w, err)
			return
		}

		go func() {
			channel <- data
		}()
	})

	c.Router.AddMatch(c.Endpoint, handler)

	return handler
}

func (c *HttpServerWorker) RemoveRoute() {
	c.Router.RemoveMatch(c.Endpoint)
}

func setupError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}
