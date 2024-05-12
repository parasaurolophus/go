// Copyright Kirk Rader 2024

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"parasaurolophus/go/ecb"
)

// Invoke ecb.Fetch manually, to support interactive debugging.
func main() {
	data, err := ecb.Fetch(ecb.HistoricalCSV, ecb.ParseCSV)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
