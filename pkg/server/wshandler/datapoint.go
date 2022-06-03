package wshandler

import (
	"context"
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/gorilla/websocket"
)

type DataPointHubClient struct {
	dataPointId shared.CommonId
	conn        *websocket.Conn
}

func NewDataPointHubClient(dataPointId shared.CommonId, conn *websocket.Conn) *DataPointHubClient {
	return &DataPointHubClient{dataPointId: dataPointId, conn: conn}
}

func (d *DataPointHubClient) Execute(ctx context.Context, message any) error {
	if series, ok := message.(shared.Series); ok {
		err := d.conn.WriteJSON(series)
		if err != nil {
			return err
		}
	}

	return errors.New("not a series")
}
