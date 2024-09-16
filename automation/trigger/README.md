Copright 2024 Kirk Rader

# Automation Trigger Events

```
package trigger // import "parasaurolophus/automation/trigger"


CONSTANTS

const (

        // E.g. turn off exterior lighting, open west-facing window coverings.
        Sunrise = Trigger("sunrise")

        // E.g. open east-facing window coverings, close west-facing ones. Note
        // that "noon" in this context refers to "solar noon," i.e. the time at
        // which the sun is at its highest altitude on any given day, not 12PM
        // local time. TODO: consider replacing this with distinct "midday" and
        // "afternoon" events but note that all of the go suncalc ports appear to
        // be based on a very out of date version of the original mourner code and
        // lack the features which would make that straightforward to implement.
        Noon = Trigger("noon")

        // E.g. turn on interior and exterior lighting, open west-facing window
        // coverings.
        Sunset = Trigger("sunset")

        // E.g. close all window coverings
        Evening = Trigger("evening")

        // E.g. set interior lights to night mode. This value is controlled by a
        // parameter passed to the NewTimer function.
        Bedtime = Trigger("bedtime")

        // E.g. turn off exterior lighting.
        Night = Trigger("night")
)

FUNCTIONS

func StartTriggersTimer(

        latitude, longitude float64,
        bedtime int,

) (

        triggers <-chan Trigger,
        terminate chan<- any,
        await <-chan any,
        err error,

)
    Launch a worker goroutine to send Trigger events at the appropriate times
    each day. It will skip events for any times-of-day that are already out of
    date when it is launched.


TYPES

type Trigger string
    Identifier for a time-of-day based automation trigger event.
```
