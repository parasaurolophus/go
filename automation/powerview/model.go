// Copyright 2024 Kirk Rader

package powerview

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
	Model map[string]Room
)
