Copyright 2024 Kirk Rader

# Hue Bridge Interface

```
package hue // import "parasaurolophus/hue"


TYPES

type HueBridge struct {
	Label   string `json:"label"`
	Address string `json:"address"`

	// Has unexported fields.
}
    Interface to the V2 API exposed by a Hue Bridge.

func New(

	label, address, key string,
	onConnect, onDisconnect func(*HueBridge),

) (

	bridge *HueBridge,
	sseData <-chan any, err error,

)
    Return a pointer to a HueBridge value initialized to communicate with the
    V2 API exposed at the specified IP address or host name, using the given
    security key. In addition, return a channel that can be used to receive SSE
    data asynchronously from the bridge.

func (hub *HueBridge) Get(resource string) (response any, err error)

func (bridge *HueBridge) Post(uri string, payload any) (any, error)
    Send a POST request to the given bridge.

func (bridge *HueBridge) Put(uri string, payload any) (any, error)
    Send a PUT request to the given bridge.
```

## Usage

```go
bridge, bridgeEvents, err := hue.New("Ground Floor", bridgeAddr, bridgeKey, nil, onDisconnect)
if err != nil {
    fmt.Fprintf(os.Stderr, "ground floor: %s\n", err.Error())
    os.Exit(2)
}

encoder := json.NewEncoder(os.Stdout)
encoder.SetIndent("", "  ")

resources, err := bridge.Get("resource")
if err != nil {
    fmt.Fprintf(os.Stderr, "ground floor: %s\n", err.Error())
}
_ = encoder.Encode(resources)

resources, err = basement.Get("resource")
if err != nil {
    fmt.Fprintf(os.Stderr, "basement: %s\n", err.Error())
}
_ = encoder.Encode(resources)

quit := make(chan any)
go func() {
    buffer := []byte{0}
    _, _ = os.Stdin.Read(buffer)
    quit <- buffer[0]
}()

encoder := json.NewEncoder(os.Stdout)
encoder.SetIndent("", "  ")

for {
    select {

    case bridgeEvent := <-bridgeEvents:
        if bridgeEvent == nil {
            return
        }
        _ = encoder.Encode(bridgeEvent)

    case <-quit:
        return
    }
}
```
