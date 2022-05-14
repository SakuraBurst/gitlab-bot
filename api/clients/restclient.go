package clients

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"
)

func init() {
	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			// UNSAFE!
			// DON'T USE IN PRODUCTION!
			InsecureSkipVerify: true,
		},
	}
}

var MockEnabled = false

type MocksTable map[string]Mock

var Mocks = MocksTable{}

type Mock struct {
	Response *http.Response
	Err      error
}

func (m MocksTable) AddMock(URL string, mock Mock) {
	m[URL] = mock
}

func (m *MocksTable) ClearMocks() {
	*m = MocksTable{}
}

func EnableMock() {
	MockEnabled = true
}

func DisableMock() {
	MockEnabled = false
}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	if MockEnabled {
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
	if MockEnabled {
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
