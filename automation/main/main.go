// Copyright 2024 Kirk Rader

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"parasaurolophus/automation/hue"
	"parasaurolophus/automation/powerview"
	"parasaurolophus/automation/trigger"
	"parasaurolophus/utilities"
	"strconv"
	"time"
)

var (
	output *os.File
)

func init() {

	var err error

	if output, err = os.Create("output.txt"); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {

	///////////////////////////////////////////////////////////////////////////
	// initialize

	groundFloorAddr, groundFloorKey, basementAddr, basementKey, powerviewAddr, latitude, longitude, ok := getEnvVars()
	if !ok {
		fmt.Fprintln(os.Stderr, "error reading environment variables")
		os.Exit(2)
	}

	quit := make(chan any)
	go func() {
		buffer := []byte{0}
		_, _ = os.Stdin.Read(buffer)
		quit <- buffer[0]
	}()

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")

	///////////////////////////////////////////////////////////////////////////
	// invoke powerview hub API

	model, err := powerview.GetModel(powerviewAddr)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}

	_ = encoder.Encode(model)

	// room := model["Default Room"]
	// scene := room.Scenes[0]
	// powerview.ActivateScene(address, scene.Id)

	///////////////////////////////////////////////////////////////////////////
	// subscribe to SSE messages from two hue briges and invoke the synchronous
	// API on each

	groundFloorEvents, groundFloorErrors, groundFloorTerminate, groundFloorAwait, err :=
		hue.SubscribeToSSE(
			groundFloorAddr,
			groundFloorKey,
			onHueConnect,
			onHueDisconnect,
		)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(4)
	}

	basementEvents, basementErrors, basementTerminate, basementAwait, err :=
		hue.SubscribeToSSE(
			basementAddr,
			basementKey,
			onHueConnect,
			onHueDisconnect,
		)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(5)
	}

	resources, err := hue.Send(groundFloorAddr, groundFloorKey, http.MethodGet, "resource", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ground floor: %s\n", err.Error())
		os.Exit(6)
	}
	_ = encoder.Encode(resources)

	resources, err = hue.Send(basementAddr, basementKey, http.MethodGet, "resource", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "basement: %s\n", err.Error())
		os.Exit(7)
	}
	_ = encoder.Encode(resources)

	///////////////////////////////////////////////////////////////////////////
	// handle the asynchronous events from all of the above, as well as
	// automation trigger events based on the given geographic coordinates and
	// bedtime hour

	err = handleEvents(

		latitude,
		longitude,
		10,

		groundFloorEvents,
		groundFloorErrors,
		groundFloorTerminate,
		groundFloorAwait,

		basementEvents,
		basementErrors,
		basementTerminate,
		basementAwait,

		quit,
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(8)
	}
}

func getEnvVars() (

	groundFloorAddr, groundFloorKey, basementAddr, basementKey, powerviewAddr string,
	latitude, longitude float64,
	ok bool,

) {

	if groundFloorAddr, ok = os.LookupEnv("GROUND_FLOOR_HUE_ADDRESS"); !ok {
		return
	}

	if groundFloorKey, ok = os.LookupEnv("GROUND_FLOOR_HUE_KEY"); !ok {
		return
	}

	if basementAddr, ok = os.LookupEnv("BASEMENT_HUE_ADDRESS"); !ok {
		return
	}

	if basementKey, ok = os.LookupEnv("BASEMENT_HUE_KEY"); !ok {
		return
	}

	if powerviewAddr, ok = os.LookupEnv("POWERVIEW_ADDRESS"); !ok {
		return
	}

	var (
		s   string
		err error
	)

	if s, ok = os.LookupEnv("LATITUDE"); !ok {
		return
	} else if latitude, err = strconv.ParseFloat(s, 64); err != nil {
		ok = false
		return
	}

	if s, ok = os.LookupEnv("LONGITUDE"); !ok {
		return
	} else if longitude, err = strconv.ParseFloat(s, 64); err != nil {
		ok = false
		return
	}

	return
}

func handleEvents(

	latitude, longitude float64,
	bedtime int,

	groundFloorEvents <-chan any,
	groundFloorErrors <-chan error,
	groundFloorTerminate chan<- any,
	groundFloorAwait <-chan any,

	basementEvents <-chan any,
	basementErrors <-chan error,
	basementTerminate chan<- any,
	basementAwait <-chan any,

	quit <-chan any,

) (

	err error,

) {

	defer utilities.CloseAndWait(groundFloorTerminate, groundFloorAwait)
	defer utilities.CloseAndWait(basementTerminate, basementAwait)

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")

	for {

		///////////////////////////////////////////////////////////////////////////
		// send automation trigger events over the course of the current day

		var (
			triggers, triggersSkipped <-chan trigger.Trigger
			triggersTerminate         chan<- any
			triggersAwait             <-chan any
		)

		if triggers, triggersSkipped, triggersTerminate, triggersAwait, err = trigger.SendTriggerEvents(latitude, longitude, bedtime); err != nil {
			return
		}

	NextDay:
		for {

			select {

			case groundFloorEvent := <-groundFloorEvents:
				_ = encoder.Encode(groundFloorEvent)

			case e := <-groundFloorErrors:
				err = e
				return

			case basementEvent := <-basementEvents:
				_ = encoder.Encode(basementEvent)

			case e := <-basementErrors:
				err = e
				return

			case event := <-triggers:
				output.WriteString(fmt.Sprintf("%s @ %s\n", event, time.Now().Format(time.DateTime)))

			case event := <-triggersSkipped:
				output.WriteString(fmt.Sprintf("skipped %s @ %s\n", event, time.Now()))

			case <-triggersAwait:
				break NextDay

			case <-quit:
				utilities.CloseAndWait(triggersTerminate, triggersAwait)
				return
			}
		}
	}
}

func onHueConnect(address string) {

	output.WriteString(fmt.Sprintf("hue hub at %s connected @ %s\n", address, time.Now().Format(time.DateTime)))
}

func onHueDisconnect(address string) {

	err := fmt.Errorf("hue hub at %s disconnected @ %s", address, time.Now().Format(time.DateTime))
	output.WriteString(err.Error())
	output.WriteString("\n")
	os.Exit(9)
}
