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

func SubscribeSSE(

        address, key string,
        onConnect, onDisconnect func(string),
        sseErrors chan<- error,

) (

        events <-chan any,
        terminate chan<- any,
        await <-chan any,

)
```
