package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadJson[T any](r *http.Request) (*T, error) {
	var response T

	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(all, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
