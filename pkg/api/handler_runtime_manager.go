package api

import (
	"net/http"
)

func GetRuntimeMangerStatusHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	s.respondJsonOk(r.Context(), w, s.RuntimeManager.GetAllDataSources())
}
