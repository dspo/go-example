package app0

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Echo interface {
	Echo() (map[string]any, error)
}

func NewEcho() Echo {
	return &echo{}
}

type echo struct {
}

func (e *echo) Echo() (map[string]any, error) {
	response, err := http.DefaultClient.Get("http://app0:80/api/v1/echo")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %v", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if ok := json.Valid(data); !ok {
		return nil, errors.Errorf("invalid json: %s", string(data))
	}

	var res = make(map[string]any)
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, errors.Errorf("failed to decode, data: %s, err: %v", string(data), err)
	}

	return res, nil
}
