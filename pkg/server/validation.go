package server

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
)

var Validate *validator.Validate = validator.New()

func ValidateFromBody[T any](r *http.Request) (*T, error) {
	var response T

	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(all, &response)
	if err != nil {
		return nil, err
	}

	err = Validate.StructCtx(r.Context(), response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
