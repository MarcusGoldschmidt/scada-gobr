package pkg

import (
	"database/sql"
	"net/http"
)

func GetDriversHandler(s *Scadagobr, w http.ResponseWriter, r *http.Request) {
	s.respondJsonOk(w, sql.Drivers())
}
