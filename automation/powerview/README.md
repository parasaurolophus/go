Copyright 2024 Kirk Rader

# PowerView Hub Wrapper

```
package powerview // import "parasaurolophus/automation/powerview"


FUNCTIONS

func ActivateScene(address string, sceneId int) (response any, err error)
    Invoke the API exposed by the PowerView hub at the specified address.


TYPES

type Model map[string]Room
    In-memory model for a powerview home.

func GetModel(address string) (model Model, err error)
    Get the in-memory representation of the current configuration for all scenes
    in all rooms from the PowerView hub at the specified address.

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
