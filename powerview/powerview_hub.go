// Copyright 2024 Kirk Rader

package powerview

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	PowerviewScene struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		RoomId int    `json:"roomId"`
	}

	PowerviewRoom struct {
		Id     int              `json:"id"`
		Name   string           `json:"name"`
		Scenes []PowerviewScene `json:"scenes,omitempty"`
	}

	PowerviewModel map[string]PowerviewRoom

	PowerviewHub struct {
		Label   string         `json:"label"`
		Address string         `json:"address"`
		Model   PowerviewModel `json:"model,omitempty"`
	}
)

type (
	scenesData struct {
		SceneData []PowerviewScene `json:"sceneData"`
	}

	roomsData struct {
		RoomData []PowerviewRoom `json:"roomData"`
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

	var scenes []PowerviewScene
	if scenes, err = powerviewHub.getScenes(); err != nil {
		return
	}

	var rooms []PowerviewRoom
	if rooms, err = powerviewHub.getRooms(scenes); err != nil {
		return
	}

	for _, room := range rooms {

		if len(room.Scenes) > 0 {
			powerviewHub.Model[room.Name] = room
		}
	}

	hub = powerviewHub
	return
}

// Send a command to activate the scene with the given id.
func (hub *PowerviewHub) ActivateScene(scene PowerviewScene) (response any, err error) {

	var resp *http.Response
	if resp, err = get(hub.Address, fmt.Sprintf("api/scenes?sceneId=%d", scene.Id)); err != nil {
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return
}

// Send a GET request to the API at the given address and URI, returning its
// response.
func get(address, uri string) (response *http.Response, err error) {

	url := fmt.Sprintf(`http://%s/%s`, address, uri)
	if response, err = http.DefaultClient.Get(url); err != nil {
		return
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		err = fmt.Errorf("%d: %s", response.StatusCode, response.Status)
		_, _ = io.ReadAll(response.Body)
		response.Body.Close()
		response = nil
	}

	return
}

func getRoomsData(address string) (response roomsData, err error) {

	var resp *http.Response
	const uri = "api/rooms"
	resp, err = get(address, uri)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return
}

func getScenesData(address string) (response scenesData, err error) {

	var resp *http.Response
	const uri = "api/scenes"
	resp, err = get(address, uri)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return
}

func (hub *PowerviewHub) getRooms(scenes []PowerviewScene) (rooms []PowerviewRoom, err error) {

	var data roomsData

	if data, err = getRoomsData(hub.Address); err != nil {
		return
	}

	r := []PowerviewRoom{}

	for _, room := range data.RoomData {

		var name []byte

		if name, err = base64.StdEncoding.DecodeString(room.Name); err != nil {
			return
		}

		s := []PowerviewScene{}
		for _, scene := range scenes {

			if scene.RoomId == room.Id {
				s = append(s, scene)
			}
		}

		r = append(r, PowerviewRoom{
			Id:     room.Id,
			Name:   string(name),
			Scenes: s,
		})
	}

	rooms = r
	return
}

func (hub *PowerviewHub) getScenes() (scenes []PowerviewScene, err error) {

	var data scenesData

	if data, err = getScenesData(hub.Address); err != nil {
		return
	}

	s := []PowerviewScene{}

	for _, scene := range data.SceneData {

		var name []byte

		if name, err = base64.StdEncoding.DecodeString(scene.Name); err != nil {
			return
		}

		s = append(s, PowerviewScene{
			Id:     scene.Id,
			Name:   string(name),
			RoomId: scene.RoomId,
		})
	}

	scenes = s
	return
}
