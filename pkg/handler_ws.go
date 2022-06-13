package pkg

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/events"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/server/wshandler"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func GetWsDataPoint(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dataPointId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		s.respondError(w, err)
		return
	}

	conn, err := wsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		s.respondError(w, err)
		return
	}

	defer conn.Close()
	client := wshandler.NewDataPointHubClient(dataPointId, conn)
	s.HubManager.AddClient(ctx, events.DataSeriesInserter+dataPointId.String(), client)
	defer s.HubManager.RemoveClient(events.DataSeriesInserter+dataPointId.String(), client)

	<-ctx.Done()
}
