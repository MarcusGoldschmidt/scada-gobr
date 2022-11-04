package pkg

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"
)

func (s *Scadagobr) respondJsonOk(ctx context.Context, w http.ResponseWriter, payload interface{}) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Ok, "")

	s.respondJson(ctx, w, http.StatusOK, payload)
}

func (s *Scadagobr) respondJsonCreated(ctx context.Context, w http.ResponseWriter, payload interface{}) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Ok, "")
	s.respondJson(ctx, w, http.StatusCreated, payload)
}

func (s *Scadagobr) respondJson(ctx context.Context, w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		s.respondError(ctx, w, err)
		return
	}
}

func (s *Scadagobr) respondBadRequest(ctx context.Context, w http.ResponseWriter, err error) {
	s.respondJson(ctx, w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func (s *Scadagobr) respondError(ctx context.Context, w http.ResponseWriter, err error) {
	span := trace.SpanFromContext(ctx)

	if _, ok := err.(*validator.InvalidValidationError); ok {
		s.respondJson(ctx, w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err, ok := err.(validator.ValidationErrors); ok {
		response := map[string]string{}

		for _, fieldError := range err {
			key := strings.ToLower(fieldError.Field()[0:1]) + fieldError.Field()[1:]
			response[key] = fieldError.Tag()
		}

		s.respondJson(ctx, w, http.StatusBadRequest, map[string]any{"errors": response})
		return
	}

	span.SetStatus(codes.Error, err.Error())
	s.Logger.Errorf("%s", err.Error())
	s.respondJson(ctx, w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
