package pkg

import (
	"net/http"
)

func GetDataSourcesRuntime(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	s.respondJsonOk(w, s.RuntimeManager.GetAllDataSources())
}
