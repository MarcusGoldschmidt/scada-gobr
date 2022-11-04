package pkg

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/shared"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

func GetDataSeriesByGroup(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	periodString := r.URL.Query()["period"]
	period := time.Minute * 60
	if len(periodString) != 0 {

		periodInt, err := strconv.ParseInt(periodString[0], 10, 64)

		if err != nil {
			s.respondError(ctx, w, err)
			return
		}

		period = time.Duration(int64(time.Minute) * periodInt)
	}

	strings := r.URL.Query()["dataPointsIds[]"]
	dataPointsIds := make([]shared.CommonId, len(strings))
	for i, value := range strings {
		id, err := uuid.Parse(value)
		if err != nil {
			s.respondError(ctx, w, err)
			return
		}

		dataPointsIds[i] = shared.CommonId(id)
	}

	now := providers.GetTimeProviderFromCtx(ctx).GetCurrentTime()

	begin := now.Add(-period)

	series, err := s.dataPointPersistence.GetPointValuesByIds(ctx, dataPointsIds, begin, now)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}

	s.respondJsonOk(ctx, w, series)
}
