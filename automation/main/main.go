package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"parasaurolophus/hue"
	"parasaurolophus/powerview"
	"strconv"
	"time"

	"github.com/sixdouglas/suncalc"
)

func main() {

	var (
		help, testHue, testPowerview bool
	)

	flag.BoolVar(&help, "help", false, "display usage and exit")
	flag.BoolVar(&testHue, "hue", false, "invoke Hue API")
	flag.BoolVar(&testPowerview, "pv", false, "invoke PowerView API")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if !(testHue || testPowerview) {
		flag.Usage()
		return
	}

	groundFloorAddr, groundFloorKey, basementAddr, basementKey, powerviewAddr, latitude, longitude, ok := getEnvVars()
	if !ok {
		fmt.Fprintln(os.Stderr, "error reading environment variables")
		os.Exit(1)
	}

	observer := suncalc.Observer{
		Latitude:  latitude,
		Longitude: longitude,
		Height:    0,
		Location:  time.Local,
	}

	times := suncalc.GetTimesWithObserver(time.Now(), observer)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(times)

	if testPowerview {
		runPowerview(powerviewAddr)
	}

	if testHue {
		runHue(groundFloorAddr, groundFloorKey, basementAddr, basementKey)
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

func handleSSE(groundFloorEvents <-chan any, basementEvents <-chan any) {

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

func onDisconnect(*hue.HueBridge) {
	panic(fmt.Errorf("hue hub disconnected"))
}

func runHue(groundFloorAddr string, groundFloorKey string, basementAddr string, basementKey string) {

	groundFloor, groundFloorEvents, err := hue.New("Ground Floor", groundFloorAddr, groundFloorKey, nil, onDisconnect)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ground floor: %s\n", err.Error())
		os.Exit(2)
	}

	basement, basementEvents, err := hue.New("Basement", basementAddr, basementKey, nil, onDisconnect)
	if err != nil {
		fmt.Fprintf(os.Stderr, "basement: %s\n", err.Error())
		os.Exit(3)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	resources, err := groundFloor.Get("resource")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ground floor: %s\n", err.Error())
	}
	_ = encoder.Encode(resources)

	resources, err = basement.Get("resource")
	if err != nil {
		fmt.Fprintf(os.Stderr, "basement: %s\n", err.Error())
	}
	_ = encoder.Encode(resources)

	handleSSE(groundFloorEvents, basementEvents)
}

func runPowerview(address string) {

	powerviewHub, err := powerview.New("Shades", address)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(4)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(powerviewHub)

	// model := powerviewHub.Model
	// room := model["Default Room"]
	// scene := room.Scenes[0]
	// powerviewHub.ActivateScene(scene)
}
