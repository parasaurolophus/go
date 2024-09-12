package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"parasaurolophus/automation"
	"parasaurolophus/automation/hue"
	"parasaurolophus/automation/powerview"
	"parasaurolophus/utilities"
	"slices"
	"strconv"
	"time"

	"github.com/sixdouglas/suncalc"
)

func main() {

	var (
		help, testHue, testPowerview, testTriggers bool
	)

	flag.BoolVar(&help, "help", false, "display usage and exit")
	flag.BoolVar(&testHue, "hue", false, "invoke Hue API")
	flag.BoolVar(&testPowerview, "pv", false, "invoke PowerView API")
	flag.BoolVar(&testTriggers, "triggers", false, "start sending automation trigger events")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if !(testHue || testPowerview || testTriggers) {
		flag.Usage()
		return
	}

	groundFloorAddr, groundFloorKey, basementAddr, basementKey, powerviewAddr, latitude, longitude, ok := getEnvVars()
	if !ok {
		fmt.Fprintln(os.Stderr, "error reading environment variables")
		os.Exit(1)
	}

	if testPowerview {
		runPowerview(powerviewAddr)
	}

	if testHue {
		runHue(groundFloorAddr, groundFloorKey, basementAddr, basementKey)
	}

	if testTriggers {
		runTriggers(latitude, longitude, 10)
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

func handleSSE(

	groundFloorEvents <-chan any,
	groundFloorTerminate chan<- any,
	groundFloorAwait <-chan any,
	basementEvents <-chan any,
	basementTerminate chan<- any,
	basementAwait <-chan any,

) {

	defer utilities.CloseAndWait(groundFloorTerminate, groundFloorAwait)
	defer utilities.CloseAndWait(basementTerminate, basementAwait)

	quit := make(chan any)
	go func() {
		buffer := []byte{0}
		_, _ = os.Stdin.Read(buffer)
		quit <- buffer[0]
	}()

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	for {

		select {

		case groundFloorEvent := <-groundFloorEvents:
			if groundFloorEvent == nil {
				return
			}
			_ = encoder.Encode(groundFloorEvent)

		case basementEvent := <-basementEvents:
			if basementEvent == nil {
				return
			}
			_ = encoder.Encode(basementEvent)

		case <-quit:
			return
		}
	}
}

func onDisconnect(address string) {

	panic(fmt.Errorf("hue hub at %s disconnected", address))
}

func runHue(groundFloorAddr string, groundFloorKey string, basementAddr string, basementKey string) {

	sseErrors := make(chan error)
	groundFloorEvents, groundFloorTerminate, groundFloorAwait := hue.SubscribeSSE(groundFloorAddr, groundFloorKey, nil, onDisconnect, sseErrors)
	basementEvents, basementTerminate, basementAwait := hue.SubscribeSSE(basementAddr, basementKey, nil, onDisconnect, sseErrors)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	resources, err := hue.Send(groundFloorAddr, groundFloorKey, http.MethodGet, "resource", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ground floor: %s\n", err.Error())
	}
	_ = encoder.Encode(resources)

	resources, err = hue.Send(basementAddr, basementKey, http.MethodGet, "resource", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "basement: %s\n", err.Error())
	}
	_ = encoder.Encode(resources)

	handleSSE(
		groundFloorEvents,
		groundFloorTerminate,
		groundFloorAwait,
		basementEvents,
		basementTerminate,
		basementAwait,
	)
}

func runPowerview(address string) {

	model, err := powerview.GetModel(address)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(4)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(model)

	// room := model["Default Room"]
	// scene := room.Scenes[0]
	// powerview.ActivateScene(address, scene.Id)
}

func runTriggers(latitude, longitude float64, bedtime int) {

	quit := make(chan any)
	go func() {
		buffer := []byte{0}
		_, _ = os.Stdin.Read(buffer)
		quit <- buffer[0]
	}()

	times := suncalc.GetTimes(time.Now(), latitude, longitude)
	display := []suncalc.DayTimeName{
		suncalc.Sunrise,
		suncalc.SolarNoon,
		suncalc.Sunset,
	}

	for k, v := range times {
		if slices.Contains(display, k) {
			fmt.Printf("%s: %s\n", v.Name, v.Value.Local())
		}
	}

	for {

		events, skipped, terminate, await, err := automation.SendTriggerEvents(latitude, longitude, bedtime)

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(5)
		}

		fmt.Println("Waiting for automation trigger events...")

	NextDay:
		for {

			select {

			case <-await:
				break NextDay

			case event := <-events:
				fmt.Printf("triggered %s @ %s\n", event, time.Now())

			case event := <-skipped:
				fmt.Printf("skipped %s @ %s\n", event, time.Now())

			case <-quit:
				close(terminate)
				<-await
				return
			}
		}
	}
}
