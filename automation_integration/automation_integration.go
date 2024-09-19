// Copyright 2024 Kirk Rader

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"parasaurolophus/automation/hue"
	"parasaurolophus/automation/powerview"
	"parasaurolophus/automation/trigger"
	"parasaurolophus/utilities"
	"strconv"
	"time"
)

func main() {

	///////////////////////////////////////////////////////////////////////////
	// initialize

	bedtime, err := parseArgs()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	groundFloorAddr, groundFloorKey, basementAddr, basementKey, powerviewAddr, latitude, longitude, ok := getEnvVars()

	if !ok {
		fmt.Fprintln(os.Stderr, "error reading environment variables")
		return
	}

	quit := make(chan any)

	go func() {
		buffer := []byte{0}
		_, _ = os.Stdin.Read(buffer)
		quit <- buffer[0]
	}()

	///////////////////////////////////////////////////////////////////////////
	// invoke powerview hub API

	powerviewHub := powerview.NewHub("Shades", powerviewAddr)
	powerviewModel, err := powerviewHub.Model()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	writeOuputJSON("powerview_model.json", powerviewModel)

	// room := model["Default Room"]
	// scene := room.Scenes[0]
	// powerview.ActivateScene(address, scene.Id)

	///////////////////////////////////////////////////////////////////////////
	// start receiving automation trigger events over the course of each day

	triggers, terminate, triggersAwait, err := trigger.StartTriggersTimer(latitude, longitude, bedtime)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	defer utilities.CloseAndWait(terminate, triggersAwait)

	///////////////////////////////////////////////////////////////////////////
	// construct the ground floor and basement hue bridge models

	groundFloorBridge := hue.NewBridge("Ground Floor", groundFloorAddr, groundFloorKey)
	groundFloorModel, err := groundFloorBridge.Model()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	writeOuputJSON("ground_floor_hue_model.json", groundFloorModel)

	basementBridge := hue.NewBridge("Basement", basementAddr, basementKey)
	basementModel, err := basementBridge.Model()

	if err != nil {
		fmt.Fprintf(os.Stderr, "basement: %s\n", err.Error())
		return
	}

	writeOuputJSON("basement_hue_model.json", basementModel)

	///////////////////////////////////////////////////////////////////////////
	// subscribe to SSE messages from both hue briges and invoke the
	// synchronous API on each

	groundFloorItems, groundFloorErrors, groundFloorTerminate, groundFloorAwait, err :=
		groundFloorBridge.Subscribe(onHueConnect, onHueDisconnect)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	defer utilities.CloseAndWait(groundFloorTerminate, groundFloorAwait)

	basementItems, basementErrors, basementTerminate, basementAwait, err :=
		basementBridge.Subscribe(onHueConnect, onHueDisconnect)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	defer utilities.CloseAndWait(basementTerminate, basementAwait)

	///////////////////////////////////////////////////////////////////////////
	// handle the asynchronous events from all of the above

	var (
		groundFloorEventsFile    *os.File
		groundFloorEventCount    = 0
		groundFloorEventsEncoder *json.Encoder
		basementEventsFile       *os.File
		basementEventCount       = 0
		basementEventsEncoder    *json.Encoder
		triggerEventsFile        *os.File
	)

	groundFloorEventsFile, err = os.Create("ground_floor_hue_events.json")

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	defer groundFloorEventsFile.Close()
	defer fmt.Fprintln(groundFloorEventsFile, "]")

	fmt.Fprintln(groundFloorEventsFile, "[")
	groundFloorEventsEncoder = json.NewEncoder(groundFloorEventsFile)
	groundFloorEventsEncoder.SetIndent("  ", "  ")
	basementEventsFile, err = os.Create("basement_hue_events.json")

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	defer basementEventsFile.Close()
	defer fmt.Fprintln(basementEventsFile, "]")

	fmt.Fprintln(basementEventsFile, "[")
	basementEventsEncoder = json.NewEncoder(basementEventsFile)
	basementEventsEncoder.SetIndent("  ", "  ")
	triggerEventsFile, err = os.Create("trigger_events.txt")

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	defer triggerEventsFile.Close()

HandleEvents:
	for {

		select {

		case groundFloorEvent := <-groundFloorItems:
			if groundFloorEventCount > 0 {
				fmt.Fprintln(groundFloorEventsFile, ",")
			}
			groundFloorEventCount++
			err = groundFloorEventsEncoder.Encode(groundFloorEvent)
			if err != nil {
				break HandleEvents
			}

		case err = <-groundFloorErrors:
			break HandleEvents

		case <-groundFloorAwait:
			break HandleEvents

		case basementEvent := <-basementItems:
			if basementEventCount > 0 {
				fmt.Fprintln(basementEventsFile, ",")
			}
			basementEventCount++
			err = basementEventsEncoder.Encode(basementEvent)
			if err != nil {
				break HandleEvents
			}

		case err = <-basementErrors:
			break HandleEvents

		case <-basementAwait:
			break HandleEvents

		case event := <-triggers:
			fmt.Fprintf(triggerEventsFile, "triggered %s @ %s\n", event, time.Now().Format(time.RFC850))

		case <-triggersAwait:
			break HandleEvents

		case <-quit:
			break HandleEvents
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func parseArgs() (

	bedtime int,
	err error,

) {

	help := false

	flagSet := flag.NewFlagSet("automation_integration", flag.ContinueOnError)
	flagSet.BoolVar(&help, "help", false, "display usage and exit")
	flagSet.IntVar(&bedtime, "bedtime", 22, "desired bedtime (0-23)")
	err = flagSet.Parse(os.Args)

	if help {
		flagSet.Usage()
		err = fmt.Errorf("exiting")
	}

	return
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

func onHueConnect(bridge hue.Bridge) {

	fmt.Printf("hue hub at %s connected @ %s\n", bridge.Label, time.Now().Format(time.RFC850))
}

func onHueDisconnect(bridge hue.Bridge) {

	err := fmt.Errorf("hue hub at %s disconnected @ %s", bridge.Label, time.Now().Format(time.RFC850))
	fmt.Fprintln(os.Stderr, err.Error())
}

func writeOuputJSON(filename string, object any) {

	file, err := os.Create(filename)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	err = encoder.Encode(object)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
