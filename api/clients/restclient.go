package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var mockEnabled = false

var Mocks map[string]response

type response struct {
	Response *http.Response
	Err      error
}

func EnableMock() {
	mockEnabled = true
}

func DisableMock() {
	mockEnabled = false
}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	if mockEnabled {
		r, ok := Mocks[url]
		if !ok {
			return nil, errors.New("Нет такого реквеста в моках")
		}
		return r.Response, r.Err
	}
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, buffer)
	request.Header = validHeaders(headers)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(request)
}

func Get(url string, headers http.Header) (*http.Response, error) {
	if mockEnabled {
		r, ok := Mocks[url]
		if !ok {
			return nil, errors.New("Нет такого реквеста в моках")
		}
		return r.Response, r.Err
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header = validHeaders(headers)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(request)
}
