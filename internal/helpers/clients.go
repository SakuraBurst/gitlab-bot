package helpers

import "net/http"

func ValidHeaders(headers http.Header) http.Header {
	if headers == nil {
		return make(http.Header)
	}
	return headers
}
