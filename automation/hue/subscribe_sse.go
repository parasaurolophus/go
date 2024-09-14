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
// Errors will be sent to the second returned channel. This function launches
// two goroutines, one of which will remain subscribed to the Hue Bridge until
// the third returned channel is closed. That worker goroutine will close the
// fourth returned channel before exiting. The other worker goroutine is
// created implicitly by calling sse.Client.SubscribeChanRaw.
func SubscribeToSSE(

	address, key string,
	onConnect, onDisconnect func(string),

) (

	events <-chan Item,
	errors <-chan error,
	terminate chan<- any,
	await <-chan any,
	err error,

) {

	// make the channels used to communicate with callers of SubscribeToSSE
	ev := make(chan Item)
	er := make(chan error)
	term := make(chan any)
	aw := make(chan any)

	// make the channel used to communicate with the worker goroutines launched
	// as a side-effect of calling SubscribeSSE
	rawEvents := make(chan *sse.Event)

	// launch a worker goroutine which consumes raw SSE messages from the hue
	// bridge at the specified address and forwards them to the events and
	// error channels, as appropriate
	go func() {

		defer close(aw)

		for {

			select {

			case <-term:
				return

			case event := <-rawEvents:
				dataReader := bytes.NewReader(event.Data)
				eventStreamReader := sse.NewEventStreamReader(dataReader, 65536)
				marshaledJSON, err := eventStreamReader.ReadEvent()
				if err != nil {
					er <- err
					continue
				}
				var datum []map[string]any
				err = json.Unmarshal(marshaledJSON, &datum)
				if err != nil {
					er <- err
					continue
				}
				walkRawMessage(datum, ev, er)
			}
		}
	}()

	// create the sse.Client used to subscribe to the raw SSE messages from the
	// bridge at the specified address
	sseClient := sse.NewClient(

		fmt.Sprintf(`https://%s/eventstream/clip/v2`, address),

		func(c *sse.Client) {

			c.Connection.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

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

	// launch a goroutine to listen for raw SSE messages, forwarding them to
	// the rawEvents channel
	if err = sseClient.SubscribeChanRaw(rawEvents); err != nil {
		utilities.CloseAndWait(term, aw)
		return
	}

	// set the returned unidirectional channels
	events = ev
	errors = er
	terminate = term
	await = aw

	return
}

// A rather crufty mechanism for handling raw SSE messages using Hue's very
// poorly designed data model.
func walkRawMessage(datum any, sseData chan<- Item, sseErrors chan<- error) {

	switch v := datum.(type) {

	case []any:
		// process each element recursively when passed a slice of any
		for _, d := range v {
			walkRawMessage(d, sseData, sseErrors)
		}

	case []map[string]any:
		// process each element recursively when passed a collection of key /
		// value pairs
		for _, d := range v {
			walkRawMessage(d, sseData, sseErrors)
		}

	case map[string]any:
		if d, ok := v["data"]; ok {
			// process the value of "data" recursively, when present
			walkRawMessage(d, sseData, sseErrors)
		} else {
			// send leaf objects to the SSE data channel
			sseData <- v
		}

	default:
		sseErrors <- fmt.Errorf("unsupported SSE payload %v of type %T", v, v)
	}
}
