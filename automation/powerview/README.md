Copyright 2024 Kirk Rader

# PowerView Hub Wrapper

```
package powerview // import "parasaurolophus/automation/powerview"


TYPES

type Hub struct {
        Rooms map[string]Room `json:"rooms"`

        // Has unexported fields.
}
    In-memory model for a powerview home.

func NewHub(address string) (hub Hub, err error)
    Get the in-memory representation of the current configuration for all scenes
    in all rooms from the PowerView hub at the specified address.

func (hub Hub) ActivateScene(scene Scene) (err error)
    Send a command to the given PowerView hub to activate the given scene.

func (hub *Hub) Refresh() (err error)
    Load the rooms data for the given hub by calling the PowerView API.

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
