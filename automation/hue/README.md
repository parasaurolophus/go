Copyright 2024 Kirk Rader

# Hue Bridge Interface

```
package hue // import "parasaurolophus/automation/hue"


FUNCTIONS

func ReceiveSSE(

        address, key string,
        onConnect, onDisconnect func(string),

) (

        items <-chan Item,
        errors <-chan error,
        terminate chan<- any,
        await <-chan any,
        err error,

)
    Start receiving SSE messages asynchronously from the Hue Bridge at the
    specified address. This function launches two goroutines as a side-effect.
    One is created implicitly by a call to sse.Client.SubscribeChanRaw.
    The other is created explicitly and can be monitored and controlled by the
    returned await and terminate channels.


TYPES

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

func (group Group) GetSceneById(id string) (scene Scene, ok bool)

type Item map[string]any
    Alias for map[string]any used as the basic data model for the Hue Bridge API
    V2.

type Model map[string]Group
    Fields of interest from /resource endpoint's response, represented in a way
    that makes them easy and efficient to use in a home automation application
    (unlike Hue's bloated and bizarre data model),

func NewModel(resources []Item) (model Model, err error)

func (model Model) GetGroupById(id string) (group Group, ok bool)

func (model Model) GetScene(groupName string, sceneName string) (scene Scene, ok bool)

type Response struct {
        Data   []Item `json:"data"`
        Errors []any  `json:"errors"`
}
    HTTP response payload structure

func SendHTTP(

        address, key, method, uri string,
        payload any,

) (

        response Response,
        err error,

)
    Invoke the V2 API exposed by the Hue Bridge at the given address.

type Scene struct {
        Name string `json:"name"`
        Id   string `json:"id"`
}
    Fields of interest from the /resource/scene/{id} endpoint's response.
```
