Copyright 2024 Kirk Rader

# Hue Bridge Interface

```
package hue // import "parasaurolophus/automation/hue"


CONSTANTS

const (

        // Keys common to all Items.
        Id   = "id"
        IdV1 = "id_v1"
        Type = "type"

        // An "owner" is present in most, but not quite all, types of Items. It is
        // used for cross-referencing items in the massive and bizarrely designed
        // data structure returned by the /resource endpoint.
        Owner = "owner"

        // Keys for the map that is the value of "owner" when present.
        Rid   = "rid"
        Rtype = "rtype"
)
    A few predefined keys for Item maps. There are countless other keys present
    in Hue Bridge payloads used for device-specific types of items and for the
    maps they contain.


FUNCTIONS

func SubscribeToSSE(

        address, key string,
        onConnect, onDisconnect func(string),

) (

        events <-chan Item,
        errors <-chan error,
        terminate chan<- any,
        await <-chan any,
        err error,

)
    Start receiving SSE messages asynchronously from the Hue Bridge at the
    specified address. SSE messages will be sent to the first returned channel.
    Errors will be sent to the second returned channel. This function launches
    two goroutines, one of which will remain subscribed to the Hue Bridge until
    the third returned channel is closed. That worker goroutine will close
    the fourth returned channel before exiting. The other worker goroutine is
    created implicitly by calling sse.Client.SubscribeChanRaw.


TYPES

type Item map[string]any
    Alias for map[string]any used as the basic data model for the Hue Bridge API
    V2.

func (item Item) Id() (id string, err error)
    Getter for item["id"].

func (item Item) IdV1() (idV1 string, err error)
    Getter for item["id_v1"].

func (item Item) Owner() (owner map[string]any, err error)
    Getter for item["owner"].

func (msg Item) OwnerRid() (rid string, err error)
    Getter for item["owner"]["rid"].

func (msg Item) OwnerType() (rtype string, err error)
    Getter for item["owner"]["rtype"].

func (item Item) Type() (typ string, err error)
    Getter for item["type"].

type Response struct {
        Data   []Item `json:"data"`
        Errors []any  `json:"errors"`
}
    HTTP response payload structure

func Send(

        address, key, method, uri string,
        payload any,

) (

        response Response,
        err error,

)
    Invoke the V2 API exposed by the Hue Bridge at the given address.
```
