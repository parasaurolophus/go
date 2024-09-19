Copyright 2024 Kirk Rader

# PowerView Hub Wrapper

```
package powerview // import "parasaurolophus/automation/powerview"


TYPES

type Hub struct {
        Label string `json:"label"`

        // Has unexported fields.
}
    In-memory model for a powerview home.

func NewHub(label, address string) Hub
    Initialize and return a Hub.

func (hub Hub) Activate(scene Scene) (err error)
    Send a command to the given PowerView hub to activate the given scene.

func (hub Hub) Model() (model Model, err error)
    Load the rooms data for the given hub by calling the PowerView API.

type Model map[string]Room
    In-memory powerview data model.

type Room struct {
        Id     int     `json:"id"`
        Name   string  `json:"name"`
        Scenes []Scene `json:"scenes,omitempty"`
}
    In-memory model for a powerview room.

type Scene struct {
        Id     int    `json:"id"`
        Name   string `json:"name"`
        RoomId int    `json:"roomId"`
}
    In-memory model for a powerview scene.
```
