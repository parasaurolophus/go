// Copyright 2024 Kirk Rader

package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Invoke the V2 API exposed by the Hue Bridge at the given address.
func SendHTTP(

	address, key, method, uri string,
	payload any,

) (

	response Response,
	err error,

) {

	url := fmt.Sprintf(`https://%s/clip/v2/%s`, address, uri)

	var body io.Reader

	if payload == nil {

		body = http.NoBody

	} else {

		var buffer []byte

		if buffer, err = json.Marshal(payload); err != nil {
			return
		}

		body = bytes.NewReader(buffer)
	}

	var req *http.Request

	if req, err = http.NewRequest(method, url, body); err != nil {
		return
	}

	req.Header.Set("hue-application-key", key)
	var resp *http.Response

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{Transport: transport}
	if resp, err = httpClient.Do(req); err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {

		err = fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
		_, _ = io.ReadAll(resp.Body)

	} else {

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&response)

	}

	return
}
