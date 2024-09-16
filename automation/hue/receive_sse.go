// Copyright 2024 Kirk Rader

package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/r3labs/sse/v2"
)

// Start receiving SSE messages asynchronously from the Hue Bridge at the
// specified address. This function launches two goroutines as a side-effect.
// One is created implicitly by a call to sse.Client.SubscribeChanRaw. The
// other is created explicitly and can be monitored and controlled by the
// returned await and terminate channels.
func ReceiveSSE(

	address, key string,
	onConnect, onDisconnect func(string),

) (

	items <-chan Item,
	errors <-chan error,
	terminate chan<- any,
	await <-chan any,
	err error,

) {

	// make the channels used to communicate with callers of this function
	ev := make(chan Item)
	er := make(chan error)
	term := make(chan any)
	aw := make(chan any)

	// set the unidirectional channels returned as values by this function
	items = ev
	errors = er
	terminate = term
	await = aw

	// make the channel used for communication between the worker goroutines
	// launched as a side-effect of calling this function
	rawEvents := make(chan *sse.Event)

	// create the sse.Client used to subscribe to the raw SSE messages from the
	// bridge at the specified address
	client := sse.NewClient(

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
	if err = client.SubscribeChanRaw(rawEvents); err != nil {
		close(aw)
		return
	}

	// launch a worker goroutine which consumes raw SSE messages from the
	// subscribed channel and forwards them to the events and error channels
	// returned by this function, as appropriate
	go func() {

		// signal that this worker goroutine has terminated
		defer close(aw)

		// signal the raw SSE listener goroutine to terminate
		defer client.Unsubscribe(rawEvents)

		for {

			select {

			// exit the worker goroutine with the terminate channel is signaled
			case <-term:
				return

			// process raw SSE messages and forward them the the items channel
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
				walkRawMessage(ev, er, datum)
			}
		}
	}()

	return
}

// A rather crufty mechanism for handling raw SSE messages in Hue's very poorly
// designed data model.
func walkRawMessage(items chan<- Item, errors chan<- error, datum any) {

	switch v := datum.(type) {

	case []any:
		// process each element recursively when passed a slice of any
		for _, d := range v {
			walkRawMessage(items, errors, d)
		}

	case []map[string]any:
		// process each element recursively when passed a collection of key /
		// value pairs
		for _, d := range v {
			walkRawMessage(items, errors, d)
		}

	case map[string]any:
		if d, ok := v["data"]; ok {
			// process the value of "data" recursively, when present
			walkRawMessage(items, errors, d)
		} else {
			// send leaf objects to the SSE data channel
			items <- v
		}

	default:
		errors <- fmt.Errorf("unsupported SSE payload %v of type %T", v, v)
	}
}
