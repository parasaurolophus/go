// Copyright Kirk Rader 2024

package utilities

import (
	"fmt"
	"io"
	"net/http"
)

// Fetch a document from the given URL.
func Fetch(url string) (readCloser io.ReadCloser, err error) {
	response, err := http.Get(url)
	defer func() {
		if err != nil && response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		err = fmt.Errorf("HTTP status %d", response.StatusCode)
		return
	}
	readCloser = response.Body
	return
}
