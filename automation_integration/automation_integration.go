// Copyright 2024 Kirk Rader

package main

import (
	"encoding/json"
	"flag"
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

	help, bedtime := parseArgs()

	if help {
		flag.Usage()
		return
	}

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
	// start receiving automation trigger events over the course of each day

	triggers, terminate, triggersAwait, err := trigger.StartTriggersTimer(latitude, longitude, bedtime)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(4)
	}

	defer utilities.CloseAndWait(terminate, triggersAwait)

	///////////////////////////////////////////////////////////////////////////
	// subscribe to SSE messages from two hue briges and invoke the synchronous
	// API on each

	groundFloorItems, groundFloorErrors, groundFloorTerminate, groundFloorAwait, err := hue.ReceiveSSE(
		groundFloorAddr,
		groundFloorKey,
		onHueConnect,
		onHueDisconnect,
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(5)
	}

	defer utilities.CloseAndWait(groundFloorTerminate, groundFloorAwait)

	basementItems, basementErrors, basementTerminate, basementAwait, err := hue.ReceiveSSE(
		basementAddr,
		basementKey,
		onHueConnect,
		onHueDisconnect,
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(6)
	}

	defer utilities.CloseAndWait(basementTerminate, basementAwait)

	resources, err := hue.SendHTTP(groundFloorAddr, groundFloorKey, http.MethodGet, "resource", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ground floor: %s\n", err.Error())
		os.Exit(7)
	}
	_ = encoder.Encode(resources)

	resources, err = hue.SendHTTP(basementAddr, basementKey, http.MethodGet, "resource", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "basement: %s\n", err.Error())
		os.Exit(8)
	}
	_ = encoder.Encode(resources)

	///////////////////////////////////////////////////////////////////////////
	// handle the asynchronous events from all of the above

	err = handleEvents(

		triggers,
		triggersAwait,
		groundFloorItems,
		groundFloorErrors,
		groundFloorAwait,
		basementItems,
		basementErrors,
		basementAwait,
		quit,
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(9)
	}
}

func parseArgs() (help bool, bedtime int) {

	flag.BoolVar(&help, "help", false, "display usage and exit")
	flag.IntVar(&bedtime, "bedtime", 22, "desired bedtime (0-23)")
	flag.Parse()
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

func handleEvents(

	triggers <-chan trigger.Trigger,
	triggersAwait <-chan any,
	groundFloorItems <-chan hue.Item,
	groundFloorErrors <-chan error,
	groundFloorAwait <-chan any,
	basementItems <-chan hue.Item,
	basementErrors <-chan error,
	basementAwait <-chan any,
	quit <-chan any,

) (

	err error,

) {

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")

	for {

		select {

		case groundFloorEvent := <-groundFloorItems:
			_ = encoder.Encode(groundFloorEvent)

		case e := <-groundFloorErrors:
			err = e
			return

		case <-groundFloorAwait:
			return

		case basementEvent := <-basementItems:
			_ = encoder.Encode(basementEvent)

		case e := <-basementErrors:
			err = e
			return

		case <-basementAwait:
			return

		case event := <-triggers:
			fmt.Fprintf(output, "triggered %s @ %s\n", event, time.Now().Format(time.RFC850))

		case <-triggersAwait:
			return

		case <-quit:
			return
		}
	}
}

func onHueConnect(address string) {

	fmt.Fprintf(output, "hue hub at %s connected @ %s\n", address, time.Now().Format(time.RFC850))
}

func onHueDisconnect(address string) {

	err := fmt.Errorf("hue hub at %s disconnected @ %s", address, time.Now().Format(time.RFC850))
	fmt.Fprintln(output, err.Error())
	os.Exit(9)
}
