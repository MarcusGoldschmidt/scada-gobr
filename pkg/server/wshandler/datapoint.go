package wshandler

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/gorilla/websocket"
	"sync"
)

type DataPointHubClient struct {
	dataPointId shared.CommonId
	conn        *websocket.Conn
	mutex       *sync.Mutex
}

func NewDataPointHubClient(dataPointId shared.CommonId, conn *websocket.Conn, mutex *sync.Mutex) *DataPointHubClient {
	return &DataPointHubClient{dataPointId: dataPointId, conn: conn, mutex: mutex}
}

func (d *DataPointHubClient) Execute(ctx context.Context, message any) error {
	if series, ok := message.(*persistence.SeriesGroupIdentifier); ok {
		d.mutex.Lock()
		defer d.mutex.Unlock()

		err := d.conn.WriteJSON(series)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("not a series")
}
