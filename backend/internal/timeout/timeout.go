package timeout

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Timeout map[string]interface{}

func (t Timeout) Value() (driver.Value, error) {
	j, err := json.Marshal(t)
	return j, err
}

func (t *Timeout) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*t, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed")
	}
	return nil
}
