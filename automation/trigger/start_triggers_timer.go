// Copyright 2024 Kirk Rader

package trigger

import (
	"fmt"
	"time"

	"github.com/sixdouglas/suncalc"
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

		fmt.Printf("sunrise: %s\n", times[suncalc.Sunrise].Value.Local().Format(time.RFC850))
		fmt.Printf("noon: %s\n", times[suncalc.SolarNoon].Value.Local().Format(time.RFC850))
		fmt.Printf("sunset: %s\n", times[suncalc.Sunset].Value.Local().Format(time.RFC850))
		fmt.Printf("evening: %s\n", times[suncalc.Night].Value.Local().Format(time.RFC850))
		fmt.Printf("bedtime: %s\n", bedtimeTime.Format(time.RFC850))
		fmt.Printf("night: %s\n", nightTime.Format(time.RFC850))

	Today:
		for {

			select {

			case <-sunriseTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.Sunrise].Value.Add(time.Minute)) {
					fmt.Printf("triggered %s @ %s\n", Sunrise, n.Format(time.RFC850))
					triggers <- Sunrise
				} else {
					fmt.Printf("skipped %s @ %s\n", Sunrise, n.Format(time.RFC850))
				}

			case <-noonTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.SolarNoon].Value.Add(time.Minute)) {
					fmt.Printf("triggered %s @ %s\n", Noon, n.Format(time.RFC850))
					triggers <- Noon
				} else {
					fmt.Printf("skipped %s @ %s\n", Noon, n.Format(time.RFC850))
				}

			case <-sunsetTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.Sunset].Value.Add(time.Minute)) {
					fmt.Printf("triggered %s @ %s\n", Sunset, n.Format(time.RFC850))
					triggers <- Sunset
				} else {
					fmt.Printf("skipped %s @ %s\n", Sunset, n.Format(time.RFC850))
				}

			case <-bedtimeTimer.C:
				n := time.Now()
				if n.Before(bedtimeTime.Add(time.Minute)) {
					fmt.Printf("triggered %s @ %s\n", Bedtime, n.Format(time.RFC850))
					triggers <- Bedtime
				} else {
					fmt.Printf("skipped %s @ %s\n", Bedtime, n.Format(time.RFC850))
				}

			case <-eveningTimer.C:
				n := time.Now()
				if n.Before(times[suncalc.Night].Value.Add(time.Minute)) {
					fmt.Printf("triggered %s @ %s\n", Evening, n.Format(time.RFC850))
					triggers <- Evening
				} else {
					fmt.Printf("skipped %s @ %s\n", Evening, n.Format(time.RFC850))
				}

			case <-nightTimer.C:
				n := time.Now()
				if n.Before(nightTime.Add(time.Minute)) {
					fmt.Printf("triggered %s @ %s\n", Night, n.Format(time.RFC850))
					triggers <- Night
				} else {
					fmt.Printf("skipped %s @ %s\n", Night, n.Format(time.RFC850))
				}
				break Today

			case <-terminate:
				fmt.Printf("automation triggers worker thread terminated @ %s\n", time.Now().Format(time.RFC850))
				return
			}
		}
	}
}
