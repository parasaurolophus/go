Copyright 2024 Kirk Rader

# Hue Bridge Interface

```
package hue // import "parasaurolophus/automation/hue"


FUNCTIONS

func Send(

        address, key, method, uri string,
        payload any,

) (

        response any,
        err error,

)
    Invoke the V2 API exposed by the Hue Bridge at the given address.

func SubscribeToSSE(

        address, key string,
        onConnect, onDisconnect func(string),

) (

        events <-chan map[string]any,
        errors <-chan error,
        terminate chan<- any,
        await <-chan any,
        err error,

)
    Start receiving SSE messages asynchronously from the Hue Bridge at the
    specified address. SSE messages will be sent to the first returned channel.
    Errors will be sent to the second returned channel. This function launches
    a goroutine which will remain subscribed to the Hue Bridge until the third
    returned channel is closed. The worker goroutine will close the fourth
    returned channel before exiting.
```
