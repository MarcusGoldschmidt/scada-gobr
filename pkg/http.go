package pkg

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
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
	if _, ok := err.(*validator.InvalidValidationError); ok {
		s.respondJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err, ok := err.(validator.ValidationErrors); ok {
		response := map[string]string{}

		for _, fieldError := range err {
			key := strings.ToLower(fieldError.Field()[0:1]) + fieldError.Field()[1:]
			response[key] = fieldError.Tag()
		}

		s.respondJson(w, http.StatusBadRequest, map[string]any{"errors": response})
		return
	}

	s.Logger.Errorf("%s", err.Error())
	s.respondJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
