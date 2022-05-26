package pkg

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(s *Scadagobr, w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		RespondError(s, w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		RespondError(s, w, err)
		return
	}
}

func RespondError(s *Scadagobr, w http.ResponseWriter, err error) {
	s.Logger.Errorf("%s", err.Error())
	RespondJSON(s, w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
