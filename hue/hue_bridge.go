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

type (
	HueBridge struct {
		label      string
		address    string
		key        string
		httpClient *http.Client
		sseClient  *sse.Client
		sseEvents  chan *sse.Event
		sseData    chan any
	}
)

func New(

	label string,
	address string,
	key string,
	onConnect,
	onDisconnect func(*HueBridge),

) (

	bridge *HueBridge,
	sseData <-chan any, err error,

) {

	bridge = &HueBridge{
		label:   label,
		address: address,
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

// Address implements automation.Hub.
func (bridge *HueBridge) Address() string {
	return bridge.address
}

func (hub *HueBridge) Get(resource string) (response any, err error) {

	u := fmt.Sprintf(`https://%s/clip/v2/%s`, hub.address, resource)
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, u, http.NoBody); err != nil {
		return
	}
	req.Header.Set("hue-application-key", hub.key)
	resp, err := hub.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
		return
	}
	decoder := json.NewDecoder(resp.Body)
	var r any
	err = decoder.Decode(&r)
	if err != nil {
		return
	}
	m, ok := r.(map[string]any)
	if !ok {
		err = fmt.Errorf("response body is of type %t", r)
		return
	}
	d, ok := m["data"]
	if !ok {
		err = fmt.Errorf("response has no data attribute")
		return
	}
	response = d
	return
}

// Return the label for the given bridge.
func (bridge *HueBridge) Label() string {
	return bridge.label
}

// Send a POST request to the given bridge.
func (bridge *HueBridge) Post(uri string, payload any) (any, error) {
	return bridge.send(http.MethodPost, uri, payload)
}

// Send a PUT request to the given bridge.
func (bridge *HueBridge) Put(uri string, payload any) (any, error) {
	return bridge.send(http.MethodPut, uri, payload)
}

func (hub *HueBridge) send(method string, uri string, payload any) (response any, err error) {
	u := fmt.Sprintf(`https://%s/clip/v2/%s`, hub.address, uri)
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	if err = encoder.Encode(payload); err != nil {
		return
	}
	body := bytes.NewReader(buffer.Bytes())
	var req *http.Request
	if req, err = http.NewRequest(method, u, body); err != nil {
		return
	}
	req.Header.Set("hue-application-key", hub.key)
	resp, err := hub.httpClient.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}
	return
}

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

func (hub *HueBridge) walkData(datum any) {

	switch v := datum.(type) {

	case []any:
		for _, d := range v {
			hub.walkData(d)
		}

	case []map[string]any:
		for _, d := range v {
			hub.walkData(d)
		}

	case map[string]any:
		if d, ok := v["data"]; ok {
			hub.walkData(d)
		} else {
			hub.sseData <- v
		}

	default:
		panic(fmt.Errorf("%v (%T)", v, v))
	}
}
