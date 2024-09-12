// Copyright 2024 Kirk Rader

package automation

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

	// E.g. turn on exterior lighting, open west-facing window coverings.
	Sunset = Trigger("sunset")

	// E.g. close all window coverings, set interior lights to night mode. This
	// value is controlled by a parameter passed to the NewTimer function.
	Bedtime = Trigger("bedtime")

	// E.g. turn off exterior lighting.
	Night = Trigger("night")
)

// Launch a worker goroutine to send Trigger events at the appropriate times on
// the current day. Triggers will be sent to the returned events channel  The
// worker goroutine will terminate after sending the "night" event or upon
// closure of the returned terminate channel. It will close the returned await
// channel before exiting. It will skip sending events for any times-of-day
// that are already more than one minute out of date when it is launched.
func SendTriggerEvents(

	latitude, longitude float64,
	bedtime int,

) (

	events <-chan Trigger,
	skipped <-chan Trigger,
	terminate chan<- any,
	await <-chan any,
	err error,

) {

	if bedtime < 0 || bedtime > 23 {
		err = fmt.Errorf("%d is not a valid bedtime", bedtime)
		return
	}

	ev := make(chan Trigger)
	events = ev

	sk := make(chan Trigger)
	skipped = sk

	aw := make(chan any)
	await = aw

	term := make(chan any)
	terminate = term

	go func() {

		defer close(aw)

		now := time.Now()
		times := suncalc.GetTimes(now, latitude, longitude)

		sunriseTimer := time.NewTimer(time.Until(times[suncalc.Sunrise].Value))
		defer sunriseTimer.Stop()

		noonTimer := time.NewTimer(time.Until(times[suncalc.SolarNoon].Value))
		defer noonTimer.Stop()

		sunsetTimer := time.NewTimer(time.Until(times[suncalc.Sunset].Value))
		defer sunsetTimer.Stop()

		bedtimeTime := time.Date(now.Year(), now.Month(), now.Day(), bedtime, 0, 0, 0, time.Local)

		if !bedtimeTime.After(times[suncalc.Sunset].Value) {
			bedtimeTime = times[suncalc.Sunset].Value.Add(30 * time.Minute)
		}

		bedtimeTimer := time.NewTimer(time.Until(bedtimeTime))
		defer bedtimeTimer.Stop()

		var nightTime time.Time

		if now.Hour() < 1 {
			nightTime = time.Date(now.Year(), now.Month(), now.Day(), 1, 1, 0, 0, time.Local)
		} else {
			tomorrow := now.Add(24 * time.Hour)
			nightTime = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 1, 1, 0, 0, time.Local)
		}

		nightTimer := time.NewTimer(time.Until(nightTime))
		defer nightTimer.Stop()

		for {

			select {

			case <-sunriseTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.Sunrise].Value.Add(time.Minute)) {
					ev <- Sunrise
				} else {
					sk <- Sunrise
				}

			case <-noonTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.SolarNoon].Value.Add(time.Minute)) {
					ev <- Noon
				} else {
					sk <- Noon
				}

			case <-sunsetTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.Sunset].Value.Add(time.Minute)) {
					ev <- Sunset
				} else {
					sk <- Sunset
				}

			case <-bedtimeTimer.C:
				n := time.Now()
				if n.Before(bedtimeTime.Add(time.Minute)) {
					ev <- Bedtime
				} else {
					sk <- Bedtime
				}

			case <-nightTimer.C:
				n := time.Now()
				if n.Before(nightTime.Add(time.Minute)) {
					ev <- Night
				} else {
					sk <- Night
				}
				return

			case <-term:
				return
			}
		}
	}()

	return
}
