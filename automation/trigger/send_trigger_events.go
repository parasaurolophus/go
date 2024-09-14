// Copyright 2024 Kirk Rader

package trigger

import (
	"fmt"
	"time"

	"github.com/sixdouglas/suncalc"
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

	go sendEvents(latitude, longitude, bedtime, ev, sk, term, aw)
	return
}

// Worker goroutine that sends Trigger events using timers based on the sun's
// position over the course of the current day.
func sendEvents(

	latitude float64,
	longitude float64,
	bedtime int,
	events chan<- Trigger,
	skipped chan<- Trigger,
	terminate <-chan any,
	await chan<- any,

) {

	defer close(await)

	now := time.Now()
	times := suncalc.GetTimes(now, latitude, longitude)

	sunriseTimer := time.NewTimer(time.Until(times[suncalc.Sunrise].Value))
	defer sunriseTimer.Stop()

	noonTimer := time.NewTimer(time.Until(times[suncalc.SolarNoon].Value))
	defer noonTimer.Stop()

	sunsetTimer := time.NewTimer(time.Until(times[suncalc.Sunset].Value))
	defer sunsetTimer.Stop()

	eveningTimer := time.NewTimer(time.Until(times[suncalc.Night].Value))
	defer eveningTimer.Stop()

	bedtimeTime := time.Date(now.Year(), now.Month(), now.Day(), bedtime, 0, 0, 0, time.Local)

	if !bedtimeTime.After(times[suncalc.Night].Value) {
		bedtimeTime = times[suncalc.Night].Value.Add(30 * time.Minute)
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

	fmt.Printf("sunrise: %s\n", times[suncalc.Sunrise].Value.Local().Format(time.RFC850))
	fmt.Printf("noon: %s\n", times[suncalc.SolarNoon].Value.Local().Format(time.RFC850))
	fmt.Printf("sunset: %s\n", times[suncalc.Sunset].Value.Local().Format(time.RFC850))
	fmt.Printf("evening: %s\n", times[suncalc.Night].Value.Local().Format(time.RFC850))
	fmt.Printf("bedtime: %s\n", bedtimeTime.Format(time.RFC850))
	fmt.Printf("night: %s\n", nightTime.Format(time.RFC850))

	for {

		select {

		case <-sunriseTimer.C:
			n := time.Now()
			if n.Before(times[suncalc.Sunrise].Value.Add(time.Minute)) {
				fmt.Printf("triggered %s @ %s\n", Sunrise, n.Format(time.RFC850))
				events <- Sunrise
			} else {
				fmt.Printf("skipped %s @ %s\n", Sunrise, n.Format(time.RFC850))
				skipped <- Sunrise
			}

		case <-noonTimer.C:
			n := time.Now()
			if n.Before(times[suncalc.SolarNoon].Value.Add(time.Minute)) {
				fmt.Printf("triggered %s @ %s\n", Noon, n.Format(time.RFC850))
				events <- Noon
			} else {
				fmt.Printf("skipped %s @ %s\n", Noon, n.Format(time.RFC850))
				skipped <- Noon
			}

		case <-sunsetTimer.C:
			n := time.Now()
			if n.Before(times[suncalc.Sunset].Value.Add(time.Minute)) {
				fmt.Printf("triggered %s @ %s\n", Sunset, n.Format(time.RFC850))
				events <- Sunset
			} else {
				fmt.Printf("skipped %s @ %s\n", Sunset, n.Format(time.RFC850))
				skipped <- Sunset
			}

		case <-bedtimeTimer.C:
			n := time.Now()
			if n.Before(bedtimeTime.Add(time.Minute)) {
				fmt.Printf("triggered %s @ %s\n", Bedtime, n.Format(time.RFC850))
				events <- Bedtime
			} else {
				fmt.Printf("skipped %s @ %s\n", Bedtime, n.Format(time.RFC850))
				skipped <- Bedtime
			}

		case <-eveningTimer.C:
			n := time.Now()
			if n.Before(times[suncalc.Night].Value.Add(time.Minute)) {
				fmt.Printf("triggered %s @ %s\n", Evening, n.Format(time.RFC850))
				events <- Evening
			} else {
				fmt.Printf("skipped %s @ %s\n", Evening, n.Format(time.RFC850))
				skipped <- Evening
			}

		case <-nightTimer.C:
			n := time.Now()
			if n.Before(nightTime.Add(time.Minute)) {
				fmt.Printf("triggered %s @ %s\n", Night, n.Format(time.RFC850))
				events <- Night
			} else {
				fmt.Printf("skipped %s @ %s\n", Night, n.Format(time.RFC850))
				skipped <- Night
			}
			fmt.Printf("automation triggers worker thread exiting after final event of the day @ %s\n", n.Format(time.RFC850))
			return

		case <-terminate:
			fmt.Printf("automation triggers worker thread terminated @ %s\n", time.Now().Format(time.RFC850))
			return
		}
	}
}
