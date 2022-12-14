package api

import (
	"database/sql"
	"net/http"
)

func GetDriversHandler(s *ScadaApi, w http.ResponseWriter, r *http.Request) {
	s.respondJsonOk(r.Context(), w, sql.Drivers())
}
