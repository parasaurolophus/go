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

	// In-memory model for a powerview scene.
	PowerviewScene struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		RoomId int    `json:"roomId"`
	}

	// In-memory model for a powerview room.
	PowerviewRoom struct {
		Id     int              `json:"id"`
		Name   string           `json:"name"`
		Scenes []PowerviewScene `json:"scenes,omitempty"`
	}

	// In-memory model for a powerview home.
	PowerviewModel map[string]PowerviewRoom

	// Interface to the API published by a powerview hub.
	PowerviewHub struct {
		Label   string         `json:"label"`
		Address string         `json:"address"`
		Model   PowerviewModel `json:"model,omitempty"`
	}
)

type (

	// Intermediate data model used to represent raw response from the
	// api/scenes endpoint.
	scenesData struct {
		SceneData []PowerviewScene `json:"sceneData"`
	}

	// Intermediate data model used to represent raw response from the
	// api/rooms endpoint.
	roomsData struct {
		RoomData []PowerviewRoom `json:"roomData"`
	}
)

// Return a pointer to a PowerviewHub with the given label, at the specified IP
// address or host name.
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

// Send a GET request to the api/rooms endpoint for the given hub, returning
// its response as a slice of PowerviewRoom values.
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

// Send a GET request to the api/scenes endpoint for the given hub, returning
// its response as a slice of PowerviewScene values.
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

// Send a GET request to the api/rooms endpoint at the specified IP address or
// host name, parsing the response as a roomsData value.
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

// Send a GET request to the api/rooms endpoint at the specified IP address or
// host name, parsing the response as a roomsData value.
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
