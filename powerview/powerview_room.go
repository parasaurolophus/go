// Copyright 2024 Kirk Rader

package powerview

import "encoding/base64"

type PowerviewRoom struct {
	Id     int              `json:"id"`
	Name   string           `json:"name"`
	Scenes []PowerviewScene `json:"scenes,omitempty"`
}

func newRoom(roomData map[string]any, scenes []PowerviewScene) (room PowerviewRoom, err error) {

	var decoded []byte
	if decoded, err = base64.StdEncoding.DecodeString(roomData["name"].(string)); err != nil {
		return
	}

	room = PowerviewRoom{
		Id:     int(roomData["id"].(float64)),
		Name:   string(decoded),
		Scenes: []PowerviewScene{},
	}

	for _, scene := range scenes {

		if scene.RoomId == room.Id {
			room.Scenes = append(room.Scenes, scene)
		}
	}

	return
}
