Copyright 2024 Kirk Rader

# PowerView Hub Interface

```
package powerview // import "parasaurolophus/powerview"


TYPES

type PowerviewHub struct {
	Label   string         `json:"label"`
	Address string         `json:"address"`
	Model   PowerviewModel `json:"model,omitempty"`
}
    Interface to the API published by a powerview hub.

func New(label, address string) (hub *PowerviewHub, err error)
    Return a pointer to a PowerviewHub with the given label, at the specified IP
    address or host name.

func (hub *PowerviewHub) ActivateScene(scene PowerviewScene) (response any, err error)
    Send a command to activate the scene with the given id.

type PowerviewModel map[string]PowerviewRoom
    In-memory model for a powerview home.

type PowerviewRoom struct {
	Id     int              `json:"id"`
	Name   string           `json:"name"`
	Scenes []PowerviewScene `json:"scenes,omitempty"`
}
    In-memory model for a powerview room.

type PowerviewScene struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	RoomId int    `json:"roomId"`
}
    In-memory model for a powerview scene.

```

## Usage

```go
powerviewHub, err := powerview.New("Shades", address)

if err != nil {
    fmt.Fprintln(os.Stderr, err.Error())
    os.Exit(4)
}

model := powerviewHub.Model
room := model["Default Room"]
scene := room.Scenes[0]
powerviewHub.ActivateScene(scene)
```
