package shared

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

func ToJson(value any) (string, error) {
	response, err := json.Marshal(value)

	return string(response), err
}

func FromJson[T any](value []byte) (*T, error) {
	var response T

	err := json.Unmarshal(value, &response)

	return &response, err
}

func FromAny[T any](data any) (*T, error) {
	var result T

	err := mapstructure.Decode(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func ValidateDataSourceType[T any](data map[string]any) ([]byte, error) {
	var result T

	err := mapstructure.Decode(data, &result)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return marshal, nil
}
