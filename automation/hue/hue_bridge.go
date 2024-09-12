package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/r3labs/sse/v2"
)

// Interface to the V2 API exposed by a Hue Bridge.
type HueBridge struct {
	Label      string          `json:"label"`
	Address    string          `json:"address"`
	key        string          `json:"-"`
	httpClient *http.Client    `json:"-"`
	sseClient  *sse.Client     `json:"-"`
	sseEvents  chan *sse.Event `json:"-"`
	sseData    chan any        `json:"-"`
}

// Return a pointer to a HueBridge value initialized to communicate with the V2
// API exposed at the specified IP address or host name, using the given
// security key. In addition, return a channel that can be used to receive SSE
// data asynchronously from the bridge.
func New(

	label, address, key string,
	onConnect, onDisconnect func(*HueBridge),

) (

	bridge *HueBridge,
	sseData <-chan any, err error,

) {

	bridge = &HueBridge{
		Label:   label,
		Address: address,
		key:     key,
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	bridge.httpClient = &http.Client{Transport: transport}
	url := fmt.Sprintf(`https://%s/eventstream/clip/v2`, address)
	bridge.sseClient = sse.NewClient(
		url,
		func(c *sse.Client) {
			c.Connection.Transport = transport
			c.Headers = map[string]string{
				"hue-application-key": key,
			}
			c.OnConnect(func(*sse.Client) {
				if onConnect != nil {
					onConnect(bridge)
				}
			})
			c.OnDisconnect(func(*sse.Client) {
				if onDisconnect != nil {
					onDisconnect(bridge)
				}
			})
		},
	)
	bridge.sseData = make(chan any)
	go bridge.handleSSE()
	bridge.sseEvents = make(chan *sse.Event)
	go bridge.subscribe()
	sseData = bridge.sseData
	return
}

func (hub *HueBridge) Get(resource string) (response any, err error) {

	u := fmt.Sprintf(`https://%s/clip/v2/%s`, hub.Address, resource)
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, u, http.NoBody); err != nil {
		return
	}

	req.Header.Set("hue-application-key", hub.key)

	var resp *http.Response
	if resp, err = hub.httpClient.Do(req); err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
		_, _ = io.ReadAll(resp.Body)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	var r any
	err = decoder.Decode(&r)
	if err != nil {
		return
	}

	var (
		m  map[string]any
		ok bool
	)

	if m, ok = r.(map[string]any); !ok {
		err = fmt.Errorf("response body is of type %t", r)
		return
	}

	var d any
	if d, ok = m["data"]; !ok {
		err = fmt.Errorf("response has no data attribute")
		return
	}

	response = d
	return
}

// Send a POST request to the given bridge.
func (bridge *HueBridge) Post(uri string, payload any) (any, error) {
	return bridge.send(http.MethodPost, uri, payload)
}

// Send a PUT request to the given bridge.
func (bridge *HueBridge) Put(uri string, payload any) (any, error) {
	return bridge.send(http.MethodPut, uri, payload)
}

// Send a request with the given payload to the V2 API exposed by the given
// bridge, returning its response payload.
func (hub *HueBridge) send(method string, uri string, payload any) (response any, err error) {

	url := fmt.Sprintf(`https://%s/clip/v2/%s`, hub.Address, uri)

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

	req.Header.Set("hue-application-key", hub.key)
	var resp *http.Response

	if resp, err = hub.httpClient.Do(req); err != nil {
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

// Goroutine launched by the HueBridge "constructor" function to handle SSE
// data asynchronously. This splits apart SSE payloads that contain multiple
// events, sending each separately to the channel returned as the second value
// from New.
func (hub *HueBridge) handleSSE() {

	for event := range hub.sseEvents {

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
		hub.walkData(datum)
	}
}

func (hub *HueBridge) subscribe() {

	if err := hub.sseClient.SubscribeChanRaw(hub.sseEvents); err != nil {
		fmt.Fprintf(os.Stderr, "hue hub subscribe: %s\n", err.Error())
	}
}

// A rather crufty mechanism for handling SSE data from Hue's very poorly
// designed data model.
func (hub *HueBridge) walkData(datum any) {

	switch v := datum.(type) {

	case []any:
		// process each element recursively when passed a slice of any
		for _, d := range v {
			hub.walkData(d)
		}

	case []map[string]any:
		// process each element recursively when passed a collection of key /
		// value pairs
		for _, d := range v {
			hub.walkData(d)
		}

	case map[string]any:
		if d, ok := v["data"]; ok {
			// process the value of "data" recursively, when present
			hub.walkData(d)
		} else {
			// send leaf objects to the SSE data channel
			hub.sseData <- v
		}

	default:
		// panic if a value of an unexpected type is encountered
		panic(fmt.Errorf("%v (%T)", v, v))
	}
}
