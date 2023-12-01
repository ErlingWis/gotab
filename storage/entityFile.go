package storage

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

type EntityData struct {
	Values       map[string]string
	ETag         int64
	LastModified time.Time
}

func (data EntityData) encode() ([]byte, error) {
	buffer := new(bytes.Buffer)

	encoder := json.NewEncoder(buffer)

	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func decodeData(io io.Reader) (EntityData, error) {
	storage := new(EntityData)

	e := json.NewDecoder(io)
	if err := e.Decode(storage); err != nil {
		return *storage, err
	}

	return *storage, nil
}