package database_helper

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB struct {
	Data interface{} `json:"data"`
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

func (j *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j.Data)
}

type JSONBArray []interface{}

func (a JSONBArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONBArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
