package pkg

import (
	"errors"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events/topics"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/models"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server/wshandler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

func GetWsTimeSeriesViewComponent(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	viewComponentId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	conn, err := wsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			s.Logger.Errorf("%s", err.Error())
		}
	}(conn)

	viewComponent, err := s.viewPersistence.GetViewComponentById(ctx, viewComponentId)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	if viewComponent.ViewType != models.TimeSeriesViewType {
		s.respondError(ctx, w, errors.New("view type is not time series"))
		return
	}

	datapointIdsInterface := viewComponent.Data["dataPointsIds"].([]interface{})

	mutex := &sync.Mutex{}

	for _, id := range datapointIdsInterface {
		parse, err := uuid.Parse(id.(string))
		if err != nil {
			s.respondError(ctx, w, err)
			return
		}
		client := wshandler.NewDataPointHubClient(parse, conn, mutex)
		s.HubManager.AddClient(ctx, topics.DataSeriesInserter+parse.String(), client)
		defer s.HubManager.RemoveClient(topics.DataSeriesInserter+parse.String(), client)
	}

	<-ctx.Done()
}
