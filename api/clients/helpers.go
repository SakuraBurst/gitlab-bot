package clients

import "net/http"

func validHeaders(headers http.Header) http.Header {
	if headers == nil {
		return make(http.Header)
	}
	return headers
}
