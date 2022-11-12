package shared

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

type ArrayJsonB[T any] []T

func (j ArrayJsonB[any]) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}

	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *ArrayJsonB[any]) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type JsonB[T any] map[string]T

func (j *JsonB[any]) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JsonB[any]) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

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
