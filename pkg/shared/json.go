package shared

import "encoding/json"

func ToJson(value any) (string, error) {
	response, err := json.Marshal(value)

	return string(response), err
}

func FromJson[T any](value []byte) (*T, error) {

	var response T

	err := json.Unmarshal(value, &response)

	return &response, err
}
