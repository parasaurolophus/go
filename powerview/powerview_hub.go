// Copyright 2024 Kirk Rader

package powerview

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	PowerviewModel map[string]PowerviewRoom

	PowerviewHub struct {
		Label   string         `json:"label"`
		Address string         `json:"address"`
		Model   PowerviewModel `json:"model,omitempty"`
	}
)

// Return a pointer to a PowerviewHub with the given label, at the specified
// address.
func New(label, address string) (hub *PowerviewHub, err error) {
	powerviewHub := &PowerviewHub{
		Label:   label,
		Address: address,
		Model:   PowerviewModel{},
	}
	var (
		roomData  []map[string]any
		sceneData []map[string]any
	)
	if roomData, err = powerviewHub.getRooms(); err != nil {
		return
	}
	if sceneData, err = powerviewHub.getScenes(); err != nil {
		return
	}
	var scenes = make([]PowerviewScene, len(sceneData))
	for i, data := range sceneData {
		if scenes[i], err = newScene(data); err != nil {
			return
		}
	}
	for _, r := range roomData {
		var room PowerviewRoom
		if room, err = newRoom(r, scenes); err != nil {
			return
		}
		if len(room.Scenes) > 0 {
			powerviewHub.Model[room.Name] = room
		}
	}
	hub = powerviewHub
	return
}

// Send a command to activate the scene with the given id.
func (hub *PowerviewHub) ActivateScene(scene PowerviewScene) (response any, err error) {

	response, err = hub.get(fmt.Sprintf("api/scenes?sceneId=%d", scene.Id))
	return
}

func (hub *PowerviewHub) get(uri string) (response any, err error) {

	url := fmt.Sprintf(`http://%s/%s`, hub.Address, uri)

	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return
}

func (hub *PowerviewHub) getData(uri, key string) (response []map[string]any, err error) {

	var r any
	r, err = hub.get(uri)
	if err != nil {
		return
	}

	m := r.(map[string]any)
	d := m[key]
	a := d.([]any)
	resp := make([]map[string]any, len(a))

	for i, e := range a {

		resp[i] = e.(map[string]any)
	}

	response = resp
	return
}

func (hub *PowerviewHub) getRooms() ([]map[string]any, error) {

	return hub.getData("api/rooms", "roomData")
}

func (hub *PowerviewHub) getScenes() ([]map[string]any, error) {

	return hub.getData("api/scenes", "sceneData")
}
