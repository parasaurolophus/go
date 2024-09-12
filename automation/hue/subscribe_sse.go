// Copyright 2024 Kirk Rader

package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/r3labs/sse/v2"
)

func SubscribeSSE(

	address, key string,
	onConnect, onDisconnect func(string),
	sseErrors chan<- error,

) (

	events <-chan any,
	terminate chan<- any,
	await <-chan any,

) {

	ev := make(chan any)
	events = ev

	term := make(chan any)
	terminate = term

	aw := make(chan any)
	await = aw

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	url := fmt.Sprintf(`https://%s/eventstream/clip/v2`, address)

	sseClient := sse.NewClient(
		url,
		func(c *sse.Client) {
			c.Connection.Transport = transport
			c.Headers = map[string]string{
				"hue-application-key": key,
			}
			c.OnConnect(func(*sse.Client) {
				if onConnect != nil {
					onConnect(address)
				}
			})
			c.OnDisconnect(func(*sse.Client) {
				if onDisconnect != nil {
					onDisconnect(address)
				}
			})
		},
	)

	sseEvents := make(chan *sse.Event)

	go func() {

		defer close(aw)

		for {

			select {

			case <-term:
				return

			case event := <-sseEvents:
				dataReader := bytes.NewReader(event.Data)
				eventStreamReader := sse.NewEventStreamReader(dataReader, 65536)
				marshaledJSON, err := eventStreamReader.ReadEvent()
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					continue
				}
				var datum []map[string]any
				err = json.Unmarshal(marshaledJSON, &datum)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					continue
				}
				walkData(datum, ev, sseErrors)
			}
		}
	}()

	go func() {
		if err := sseClient.SubscribeChanRaw(sseEvents); err != nil {
			sseErrors <- err
		}
	}()

	return
}

// A rather crufty mechanism for handling SSE data from Hue's very poorly
// designed data model.
func walkData(datum any, sseData chan<- any, sseErrors chan<- error) {

	switch v := datum.(type) {

	case []any:
		// process each element recursively when passed a slice of any
		for _, d := range v {
			walkData(d, sseData, sseErrors)
		}

	case []map[string]any:
		// process each element recursively when passed a collection of key /
		// value pairs
		for _, d := range v {
			walkData(d, sseData, sseErrors)
		}

	case map[string]any:
		if d, ok := v["data"]; ok {
			// process the value of "data" recursively, when present
			walkData(d, sseData, sseErrors)
		} else {
			// send leaf objects to the SSE data channel
			sseData <- v
		}

	default:
		sseErrors <- fmt.Errorf("%v (%T)", v, v)
	}
}
