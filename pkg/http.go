package pkg

import (
	"encoding/json"
	"net/http"
)

func (s *Scadagobr) respondJsonOk(w http.ResponseWriter, payload interface{}) {
	s.respondJson(w, http.StatusOK, payload)
}

func (s *Scadagobr) respondJsonCreated(w http.ResponseWriter, payload interface{}) {
	s.respondJson(w, http.StatusCreated, payload)
}

func (s *Scadagobr) respondJson(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		s.respondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		s.respondError(w, err)
		return
	}
}

func (s *Scadagobr) respondError(w http.ResponseWriter, err error) {
	s.Logger.Errorf("%s", err.Error())
	s.respondJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
