package httputil

import "encoding/json"

type HTTPError struct {
	// HTTP status code to write into HTTP header.
	// This field should not be marshaled into response
	// JSON body.
	StatusCode int `json:"-"`
	// Deepshare specific error code.
	// Error code can be pragmatically consumed.
	Code int `json:"code"`
	// The error message of the error code.
	// Error message can be printed out and consumed by human.
	Message string `json:"message"`
}

func (he HTTPError) Error() string {
	b, err := json.Marshal(&he)
	if err != nil {
		panic("unexpected json marshal error")
	}
	return string(b)
}
