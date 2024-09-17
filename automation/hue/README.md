Copyright 2024 Kirk Rader

# Hue Bridge Interface

```
package hue // import "parasaurolophus/automation/hue"


TYPES

type Bridge struct {
        Label  string           `json:"label"`
        Groups map[string]Group `json:"groups"`

        // Has unexported fields.
}
    Fields of interest from /resource endpoint's response, represented in a way
    that makes them easy and efficient to use in a home automation application
    (unlike hue's bloated and bizarre data model),

func NewBridge(label, address, key string) (model Bridge, err error)
    Load a Bridge from the specified hue bridge.

func (model Bridge) ActivateScene(groupName, sceneName string) (err error)
    Send a PUT command to activate the specified scene.

func (model Bridge) ReceiveSSE(

        onConnect, onDisconnect func(string),

) (

        items <-chan Item,
        errors <-chan error,
        terminate chan<- any,
        await <-chan any,
        err error,

)

func (model *Bridge) Refresh() (err error)
    Send a GET command to update the given Model.

func (model Bridge) Send(method, uri string, payload any) (response Response, err error)

func (model Bridge) SetGroupState(groupName string, on bool) (err error)
    Send a PUT command to turn on or off the specified group.

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
