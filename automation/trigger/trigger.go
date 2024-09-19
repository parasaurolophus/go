// Copyright 2024 Kirk Rader

package trigger

import (
	"fmt"
	"time"

	"github.com/sixdouglas/suncalc"
)

// Identifier for a time-of-day based automation trigger event.
type Trigger string

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

// Launch a worker goroutine to send Trigger events at the appropriate times
// each day. It will skip events for any times-of-day that are already out of
// date when it is launched.
func StartTriggersTimer(

	latitude, longitude float64,
	bedtime int,

) (

	triggers <-chan Trigger,
	terminate chan<- any,
	await <-chan any,
	err error,

) {

	if bedtime < 0 || bedtime > 23 {
		err = fmt.Errorf("%d is not a valid bedtime", bedtime)
		return
	}

	// make the channels used to send and receive messages and signals
	trig := make(chan Trigger)
	term := make(chan any)
	aw := make(chan any)

	// set the unidirectional channels exposed to users of triggersTimer
	triggers = trig
	terminate = term
	await = aw

	// start the worker goroutine
	go worker(latitude, longitude, bedtime, trig, term, aw)

	return
}

func worker(

	latitude, longitude float64,
	bedtime int,
	triggers chan Trigger,
	terminate chan any,
	await chan any,

) {

	defer close(await)

	for {

		now := time.Now()
		base := time.Date(now.Year(), now.Month(), now.Day(), 1, 1, 0, 0, time.Local)
		times := suncalc.GetTimes(base, latitude, longitude)

		bedtimeTime := time.Date(now.Year(), now.Month(), now.Day(), bedtime, 0, 0, 0, time.Local)

		if !bedtimeTime.After(times[suncalc.Night].Value) {
			bedtimeTime = times[suncalc.Night].Value.Add(30 * time.Minute)
		}

		var nightTime time.Time

		if now.Hour() < 1 {
			nightTime = time.Date(now.Year(), now.Month(), now.Day(), 1, 1, 0, 0, time.Local)
		} else {
			tomorrow := now.Add(24 * time.Hour)
			nightTime = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 1, 1, 0, 0, time.Local)
		}

		sunriseTimer := time.NewTimer(time.Until(times[suncalc.Sunrise].Value))
		noonTimer := time.NewTimer(time.Until(times[suncalc.SolarNoon].Value))
		sunsetTimer := time.NewTimer(time.Until(times[suncalc.Sunset].Value))
		eveningTimer := time.NewTimer(time.Until(times[suncalc.Night].Value))
		bedtimeTimer := time.NewTimer(time.Until(bedtimeTime))
		nightTimer := time.NewTimer(time.Until(nightTime))

	Today:
		for {

			select {

			case <-sunriseTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.Sunrise].Value.Add(time.Minute)) {
					triggers <- Sunrise
				}

			case <-noonTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.SolarNoon].Value.Add(time.Minute)) {
					triggers <- Noon
				}

			case <-sunsetTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.Sunset].Value.Add(time.Minute)) {
					triggers <- Sunset
				}

			case <-bedtimeTimer.C:
				n := time.Now()
				if n.Before(bedtimeTime.Add(time.Minute)) {
					triggers <- Bedtime
				}

			case <-eveningTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.Night].Value.Add(time.Minute)) {
					triggers <- Evening
				}

			case <-nightTimer.C:
				n := time.Now()
				if n.Before(nightTime.Add(time.Minute)) {
					triggers <- Night
				}
				break Today

			case <-terminate:
				return
			}
		}
	}
}
