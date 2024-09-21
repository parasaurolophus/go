// Copyright 2024 Kirk Rader

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"parasaurolophus/automation/hue"
)

func main() {

	var (
		help   bool
		uri    string
		usePut bool
	)

	flag.BoolVar(&help, "help", false, "display usage and exit")
	flag.StringVar(&uri, "uri", "resource", "Hue Bridge API v2 URL suffix")
	flag.BoolVar(&usePut, "put", false, "use PUT rather than GET (with JSON payload from stdin)")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	var (
		address, key string
		ok           bool
		err          error
	)

	if address, ok = os.LookupEnv("BASEMENT_HUE_ADDRESS"); !ok {
		fmt.Fprintln(os.Stderr, `$BASEMENT_HUE_ADDRESS not defined`)
		return
	}

	if key, ok = os.LookupEnv("BASEMENT_HUE_KEY"); !ok {
		fmt.Fprintln(os.Stderr, `$BASEMENT_HUE_KEY not defined`)
		return
	}

	var payload any

	if usePut {

		decoder := json.NewDecoder(os.Stdin)
		if err := decoder.Decode(&payload); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}

	var (
		bridge   = hue.NewBridge("Basement", address, key)
		response hue.Response
	)

	if response, err = bridge.Send(http.MethodGet, uri, nil); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(response); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
