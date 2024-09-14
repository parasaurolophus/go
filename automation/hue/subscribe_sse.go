// Copyright 2024 Kirk Rader

package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"parasaurolophus/utilities"

	"github.com/r3labs/sse/v2"
)

// Start receiving SSE messages asynchronously from the Hue Bridge at the
// specified address. SSE messages will be sent to the first returned channel.
// Errors will be sent to the second returned channel. This function launches a
// goroutine which will remain subscribed to the Hue Bridge until the third
// returned channel is closed. The worker goroutine will close the fourth
// returned channel before exiting.
func SubscribeToSSE(

	address, key string,
	onConnect, onDisconnect func(string),

) (

	events <-chan map[string]any,
	errors <-chan error,
	terminate chan<- any,
	await <-chan any,
	err error,

) {

	ev := make(chan map[string]any)
	events = ev

	er := make(chan error)
	errors = er

	term := make(chan any)
	terminate = term

	aw := make(chan any)
	await = aw

	sseEvents := make(chan *sse.Event)
	go process(ev, term, aw, er, sseEvents)
	err = subscribe(address, key, onConnect, onDisconnect, sseEvents)

	if err != nil {
		utilities.CloseAndWait(term, aw)
	}

	return
}

// Worker goroutine that consumes raw SSE messages and forwards them to the
// output events channel.
func process(

	events chan<- map[string]any,
	terminate <-chan any,
	await chan<- any,
	sseErrors chan<- error,
	sseEvents <-chan *sse.Event,

) {

	defer close(await)

	for {

		select {

		case <-terminate:
			return

		case event := <-sseEvents:
			dataReader := bytes.NewReader(event.Data)
			eventStreamReader := sse.NewEventStreamReader(dataReader, 65536)
			marshaledJSON, err := eventStreamReader.ReadEvent()
			if err != nil {
				sseErrors <- err
				continue
			}
			var datum []map[string]any
			err = json.Unmarshal(marshaledJSON, &datum)
			if err != nil {
				sseErrors <- err
				continue
			}
			walkData(datum, events, sseErrors)
		}
	}
}

// Worker go routine that subscribes to raw SSE messages from the specified Hue
// Bridge.
func subscribe(

	address, key string,
	onConnect, onDisconnect func(string),
	sseEvents chan *sse.Event,

) (

	err error,

) {

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

	err = sseClient.SubscribeChanRaw(sseEvents)
	return
}

// A rather crufty mechanism for handling SSE data from Hue's very poorly
// designed data model.
func walkData(datum any, sseData chan<- map[string]any, sseErrors chan<- error) {

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
		sseErrors <- fmt.Errorf("unsupported SSE payload %v of type %T", v, v)
	}
}
