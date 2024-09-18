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
	Scene struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		RoomId int    `json:"roomId"`
	}

	// In-memory model for a powerview room.
	Room struct {
		Id     int     `json:"id"`
		Name   string  `json:"name"`
		Scenes []Scene `json:"scenes,omitempty"`
	}

	// In-memory model for a powerview home.
	Hub struct {
		Rooms map[string]Room `json:"rooms"`

		address string
	}
)

type (

	// Intermediate data model used to represent a raw response from the
	// api/scenes endpoint.
	scenesData struct {
		SceneData []Scene `json:"sceneData"`
	}

	// Intermediate data model used to represent a raw response from the
	// api/rooms endpoint.
	roomsData struct {
		RoomData []Room `json:"roomData"`
	}

	// Type constraint used by generic powerview.getData function.
	powerviewData interface {
		roomsData | scenesData
	}
)

// Get the in-memory representation of the current configuration for all scenes
// in all rooms from the PowerView hub at the specified address.
func NewHub(address string) (hub Hub, err error) {

	hub = Hub{
		address: address,
		Rooms:   map[string]Room{},
	}

	if err = hub.Refresh(); err != nil {
		return
	}

	return
}

// Load the rooms data for the given hub by calling the PowerView API.
func (hub *Hub) Refresh() (err error) {

	var scenes []Scene

	if scenes, err = getScenes(hub.address); err != nil {
		return
	}

	var rooms []Room

	if rooms, err = getRooms(hub.address, scenes); err != nil {
		return
	}

	for _, room := range rooms {

		if len(room.Scenes) > 0 {
			hub.Rooms[room.Name] = room
		}
	}

	return
}

// Send a command to the given PowerView hub to activate the given scene.
func (hub Hub) ActivateScene(scene Scene) (err error) {

	url := fmt.Sprintf(`http://%s/scenes?sceneId=%d`, hub.address, scene.Id)

	var resp *http.Response
	if resp, err = http.DefaultClient.Get(url); err != nil {
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		return
	}

	var response any
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return
}

// Invoke the API exposed by the PowerView hub at the specified address.
func getData[Value powerviewData](address, uri string) (response Value, err error) {

	url := fmt.Sprintf(`http://%s/%s`, address, uri)

	var resp *http.Response
	if resp, err = http.DefaultClient.Get(url); err != nil {
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return
}

// Send a GET request to the api/rooms endpoint for the given hub, returning
// its response as a slice of PowerviewRoom values.
func getRooms(address string, scenes []Scene) (rooms []Room, err error) {

	var data roomsData

	if data, err = getData[roomsData](address, "api/rooms"); err != nil {
		return
	}

	r := []Room{}

	for _, room := range data.RoomData {

		var name []byte

		if name, err = base64.StdEncoding.DecodeString(room.Name); err != nil {
			return
		}

		s := []Scene{}
		for _, scene := range scenes {

			if scene.RoomId == room.Id {
				s = append(s, scene)
			}
		}

		r = append(r, Room{
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
func getScenes(address string) (scenes []Scene, err error) {

	var data scenesData

	if data, err = getData[scenesData](address, "api/scenes"); err != nil {
		return
	}

	s := []Scene{}

	for _, scene := range data.SceneData {

		var name []byte

		if name, err = base64.StdEncoding.DecodeString(scene.Name); err != nil {
			return
		}

		s = append(s, Scene{
			Id:     scene.Id,
			Name:   string(name),
			RoomId: scene.RoomId,
		})
	}

	scenes = s
	return
}
