// Copyright 2024 Kirk Rader

package powerview

import "encoding/base64"

type PowerviewScene struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	RoomId int    `json:"room_id"`
}

func newScene(data map[string]any) (scene PowerviewScene, err error) {

	var decoded []byte
	if decoded, err = base64.StdEncoding.DecodeString(data["name"].(string)); err != nil {
		return
	}

	scene = PowerviewScene{
		Id:     int(data["id"].(float64)),
		Name:   string(decoded),
		RoomId: int(data["roomId"].(float64)),
	}

	return
}
