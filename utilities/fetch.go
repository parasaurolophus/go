package utilities

import (
	"fmt"
	"io"
	"net/http"
)

// Fetch a document from the given URL.
func Fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}
	return resp.Body, nil
}
