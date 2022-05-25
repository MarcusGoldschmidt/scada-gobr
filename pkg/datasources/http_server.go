package datasources

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"scadagobr/pkg/persistence"
	"scadagobr/pkg/server"
	"scadagobr/pkg/shared"
	"strings"
	"sync/atomic"
	"time"
)

type HttpServerDataPoint struct {
	*DataPointCommon
	rowIdentifier string
	dateFormat    string
}

func (r HttpServerDataPoint) Id() shared.CommonId {
	return r.id
}

func (r HttpServerDataPoint) Name() string {
	return r.name
}

type HttpServerDataSource struct {
	*DataSourceCommon
	period     time.Duration
	dataPoints []*HttpServerDataPoint

	router       *server.Router
	endpoint     string
	user         string
	passwordHash string
}

type HttpServerDataSourceRuntime struct {
	id          shared.CommonId
	dataSource  HttpServerDataSource
	persistence persistence.DataPointPersistence
	atmDone     int32
}

func (r *HttpServerDataSource) AddDataPoint(dp *HttpServerDataPoint) {
	r.dataPoints = append(r.dataPoints, dp)
}

func (r HttpServerDataSource) Id() shared.CommonId {
	return r.id
}

func (r HttpServerDataSource) Name() string {
	return r.name
}

func (r HttpServerDataSource) IsEnable() bool {
	return r.isEnable
}

func (r HttpServerDataSource) GetDataPoints() []Datapoint {
	datapoint := make([]Datapoint, len(r.dataPoints))
	for i, v := range r.dataPoints {
		datapoint[i] = Datapoint(v)
	}
	return datapoint
}

func (r HttpServerDataSource) CreateRuntime(ctx context.Context, p persistence.DataPointPersistence) (DataSourceRuntime, error) {
	rt := HttpServerDataSourceRuntime{uuid.New(), r, p, 0}
	return rt, nil
}

func (c HttpServerDataSourceRuntime) GetDataSource() DataSource {
	return c.dataSource
}

type request struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func (c HttpServerDataSourceRuntime) Run(ctx context.Context, shutdownCompleteChan chan shared.CommonId) error {
	defer func() {
		shutdownCompleteChan <- c.id
	}()

	// TODO: verify best number
	channel := make(chan []*request, 1024)

	c.GerOrAddRoute(channel)
	defer c.RemoveRoute()
	defer close(channel)
	defer atomic.AddInt32(&c.atmDone, 1)

	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-channel:

			dict := make(map[string]*shared.CommonId)

			for _, point := range c.dataSource.dataPoints {
				dict[point.rowIdentifier] = &point.id
			}

			series := make([]*shared.IdSeries, 0)
			for _, d := range data {
				if id, ok := dict[d.Name]; ok {
					series = append(series, shared.NewIdSeries(*id, d.Value, d.Timestamp))
				}
			}

			err := c.persistence.AddDataPointValues(ctx, series)
			if err != nil {
				return err
			}
		}
	}
}

func (c *HttpServerDataSourceRuntime) GerOrAddRoute(channel chan []*request) http.Handler {
	if data, ok := c.dataSource.router.VerifyMatch(c.dataSource.endpoint); ok {
		return data
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		basicAuth := r.Header["Authorization"][0][6:]

		split := strings.Split(basicAuth, ":")

		hash := sha256.Sum256([]byte(split[1]))

		if split[0] != c.dataSource.user && c.dataSource.passwordHash != string(hash[:]) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if c.atmDone == 1 {
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

	c.dataSource.router.AddMatch(c.dataSource.endpoint, handler)

	response := http.Handler(handler)

	return response
}

func setupError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

func (c *HttpServerDataSourceRuntime) RemoveRoute() {
	c.dataSource.router.RemoveMatch(c.dataSource.endpoint)
}
