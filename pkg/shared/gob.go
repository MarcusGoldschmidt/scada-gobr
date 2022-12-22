package shared

import (
	"bytes"
	"encoding/gob"
)

func ToBytesGob(data any) ([]byte, error) {
	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)

	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func FromBytesGob[T any](data []byte) (*T, error) {
	var result T

	buf := bytes.NewBuffer(data)

	err := gob.NewDecoder(buf).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
