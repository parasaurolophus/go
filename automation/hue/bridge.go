// Copyright 2024 Kirk Rader

package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"parasaurolophus/utilities"

	"github.com/r3labs/sse/v2"
)

type (

	// Fields of interest from the /resource/scene/{id} endpoint's response.
	Scene struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}

	// Fields of interest from /resource/bridge_home, /resource/room/{id} or
	// /resource/zone/{id} endpoints' responses, plus relevant fields from
	// related scene and grouped_light resources for the given group.
	Group struct {
		Name           string           `json:"name"`
		Id             string           `json:"id"`
		Type           string           `json:"type"`
		On             bool             `json:"on"`
		GroupedLightId string           `json:"grouped_light_id"`
		Scenes         map[string]Scene `json:"scenes,omitempty"`
	}

	// Fields of interest from /resource endpoint's response, represented in a
	// way that makes them easy and efficient to use in a home automation
	// application (unlike hue's bloated and bizarre data model),
	Bridge struct {
		Label  string           `json:"label"`
		Groups map[string]Group `json:"groups"`

		address string
		key     string
	}
)

// Load a Bridge from the specified hue bridge.
func NewBridge(label, address, key string) (model Bridge, err error) {

	model = Bridge{
		Label:   label,
		address: address,
		key:     key,
	}

	err = model.Refresh()

	return
}

// Send a PUT command to activate the specified scene.
func (model Bridge) ActivateScene(groupName, sceneName string) (err error) {

	var (
		group Group
		scene Scene
		ok    bool
	)

	if group, ok = model.Groups[groupName]; !ok {
		err = fmt.Errorf("no group %s found in bridge %s", groupName, model.Label)
		return
	}

	if scene, ok = group.Scenes[sceneName]; !ok {
		err = fmt.Errorf("no scene %s found in group %s", sceneName, group.Name)
		return
	}

	uri := fmt.Sprintf("/resource/scene/%s", scene.Id)
	payload := map[string]any{"recall": map[string]any{"action": "active"}}
	_, err = model.Send(http.MethodPut, uri, payload)
	return
}

func (model Bridge) ReceiveSSE(

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

		fmt.Sprintf(`https://%s/eventstream/clip/v2`, model.address),

		func(c *sse.Client) {

			c.Connection.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			c.Headers = map[string]string{
				"hue-application-key": model.key,
			}

			c.OnConnect(func(*sse.Client) {
				if onConnect != nil {
					onConnect(model.Label)
				}
			})

			c.OnDisconnect(func(*sse.Client) {
				if onDisconnect != nil {
					onDisconnect(model.Label)
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

// Send a GET command to update the given Model.
func (model *Bridge) Refresh() (err error) {

	var response Response

	model.Groups = map[string]Group{}

	if response, err = model.Send(http.MethodGet, "/resource", nil); err != nil {
		return
	}

	for _, resource := range response.Data {

		var groupType string

		if groupType, err = utilities.GetJSONPath[string](resource, "type"); err != nil {
			return
		}

		if !(groupType == "bridge_home" || groupType == "room" || groupType == "zone") {
			continue
		}

		var groupName string

		if groupName, _ = utilities.GetJSONPath[string](resource, "metadata", "name"); groupName == "" {
			groupName = "All Lights"
		}

		var groupId string

		if groupId, err = utilities.GetJSONPath[string](resource, "id"); err != nil {
			return
		}

		var (
			group Group
			ok    bool
		)

		if group, ok = model.Groups[groupName]; !ok {

			var (
				resourceType, ownerId, resourceId string
				groupedLightState                 bool
			)

			for _, r := range response.Data {

				if resourceType, err = utilities.GetJSONPath[string](r, "type"); err != nil {
					return
				}

				if resourceType != "grouped_light" {
					continue
				}

				if ownerId, err = utilities.GetJSONPath[string](r, "owner", "rid"); err != nil {
					return
				}

				if ownerId != groupId {
					continue
				}

				if resourceId, err = utilities.GetJSONPath[string](r, "id"); err != nil {
					return
				}

				if groupedLightState, err = utilities.GetJSONPath[bool](r, "on", "on"); err != nil {
					return
				}

				break
			}

			group = Group{
				Name:           groupName,
				Id:             groupId,
				Type:           groupType,
				GroupedLightId: resourceId,
				On:             groupedLightState,
				Scenes:         map[string]Scene{},
			}

			for _, r := range response.Data {

				if resourceType, err = utilities.GetJSONPath[string](r, "type"); err != nil {
					return
				}

				if resourceType != "scene" {
					continue
				}

				if ownerId, err = utilities.GetJSONPath[string](r, "group", "rid"); err != nil {
					return
				}

				if ownerId != groupId {
					continue
				}

				var sceneName, sceneId string

				if sceneName, err = utilities.GetJSONPath[string](r, "metadata", "name"); err != nil {
					return
				}

				if sceneId, err = utilities.GetJSONPath[string](r, "id"); err != nil {
					return
				}

				scene := Scene{
					Name: sceneName,
					Id:   sceneId,
				}

				group.Scenes[sceneName] = scene
			}
		}

		model.Groups[groupName] = group
	}

	return
}

func (model Bridge) Send(method, uri string, payload any) (response Response, err error) {

	url := fmt.Sprintf(`https://%s/clip/v2/%s`, model.address, uri)

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

	req.Header.Set("hue-application-key", model.key)
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

// Send a PUT command to turn on or off the specified group.
func (model Bridge) SetGroupState(groupName string, on bool) (err error) {

	var (
		group Group
		ok    bool
	)

	if group, ok = model.Groups[groupName]; !ok {
		err = fmt.Errorf("no group %s found in bridge %s", groupName, model.Label)
		return
	}

	uri := fmt.Sprintf("/resource/grouped_light/%s", group.GroupedLightId)

	payload := map[string]any{
		"on": map[string]any{
			"on": on,
		},
	}

	_, err = model.Send(http.MethodPut, uri, payload)
	return
}

// A rather crufty mechanism for handling raw SSE messages in hue's very poorly
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
