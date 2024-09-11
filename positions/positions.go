// Copyright Kirk Rader 2024

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sixdouglas/suncalc"
)

func main() {

	var (
		latitude  float64
		longitude float64
		err       error
	)

	if s, ok := os.LookupEnv("LATITUDE"); !ok {
		fmt.Fprintf(os.Stderr, "$LATITUDE undefined\n")
		os.Exit(1)
	} else if latitude, err = strconv.ParseFloat(s, 64); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	if s, ok := os.LookupEnv("LONGITUDE"); !ok {
		fmt.Fprintf(os.Stderr, "$LONGITUDE undefined\n")
		os.Exit(3)
	} else if longitude, err = strconv.ParseFloat(s, 64); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(4)
	}

	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end := t.Add(time.Hour*24 - time.Second)

	fmt.Println("time,altitude,azimuth")

	for {

		if t.After(end) {
			break
		}

		position := suncalc.GetPosition(t, latitude, longitude)
		fmt.Printf("%s,%f,%f\n", t.Format(time.TimeOnly), position.Altitude, position.Azimuth)
		t = t.Add(time.Minute)
	}
}
