Copyright 2024 Kirk Rader

# Hue Bridge Interface

```
package hue // import "parasaurolophus/automation/hue"


TYPES

type Bridge struct {
        Label string `json:"label"`

        // Has unexported fields.
}
    Parameters required to access the API exposed by a particular Hue Bridge.

func NewBridge(label, address, key string) Bridge
    Initialize and return a Bridge.

func (bridge Bridge) Activate(scene Scene) (err error)
    Send a PUT command to activate the given scene.

func (bridge Bridge) Model() (groups Model, err error)
    Send a GET command to return the Model representing the current state of the
    given Model.

func (bridge Bridge) Put(group Group) (err error)
    Send a PUT command to turn on or off the specified group.

func (bridge Bridge) Send(method, uri string, payload any) (response Response, err error)

func (bridge Bridge) Subscribe(

        onConnect, onDisconnect func(Bridge),

) (

        items <-chan Item,
        errors <-chan error,
        terminate chan<- any,
        await <-chan any,
        err error,

)
    Subscribe to SSE messages from the given Bridge.

type Group struct {
        Name           string           `json:"name"`
        Id             string           `json:"id"`
        Type           string           `json:"type"`
        On             bool             `json:"on"`
        GroupedLightId string           `json:"grouped_light_id"`
        Scenes         map[string]Scene `json:"scenes,omitempty"`
}
    Fields of interest from /resource/bridge_home, /resource/room/{id} or
    /resource/zone/{id} endpoints' responses, plus relevant fields from related
    scene and grouped_light resources for the given group.

type Item map[string]any
    Alias for map[string]any used as the basic data model for the Hue Bridge API
    V2.

type Model map[string]Group
    Fields of interest from the Hue API V2 data model, transformed into a
    useable structure (which Hue's bizzare and over-engineered structure is
    not).

type Response struct {
        Data   []Item `json:"data"`
        Errors []any  `json:"errors"`
}
    HTTP response payload structure

type Scene struct {
        Name string `json:"name"`
        Id   string `json:"id"`
}
    Fields of interest from the /resource/scene/{id} endpoint's response.
```
