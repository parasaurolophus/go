// Copyright Kirk Rader 2024

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"parasaurolophus/go/ecb"
	"parasaurolophus/go/utilities"
)

// Invoke ecb.Fetch manually, to support interactive debugging.
func main() {
	format := flag.String("format", "csv", "csv or xml")
	version := flag.String("version", "daily", "daily, ninety or historical")
	parse := flag.Bool("parse", false, "true or false")
	flag.Parse()
	var url string
	var parser ecb.Parser
	if *format == "csv" {
		parser = ecb.ParseCSV
		switch *version {
		case "daily":
			url = ecb.DAILY_CSV_URL
		case "historical":
			url = ecb.HISTORICAL_CSV_URL
		default:
			fmt.Fprintf(os.Stderr, `"%s" is not a valid version for csv\n`, *version)
			flag.Usage()
			os.Exit(1)
		}
	} else if *format == "xml" {
		parser = ecb.ParseXML
		switch *version {
		case "daily":
			url = ecb.DAILY_XML_URL
		case "historical":
			url = ecb.HISTORICAL_XML_URL
		case "ninety":
			url = ecb.NINETY_DAY_XML_URL
		default:
			fmt.Fprintf(os.Stderr, `"%s" is not a valid version for csv\n`, *version)
			flag.Usage()
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, `"%s" is not a valid format\n`, *format)
		flag.Usage()
		os.Exit(2)
	}
	source, err := utilities.Fetch(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
	var documents []io.ReadCloser
	if *format == "csv" {
		documents, err = utilities.Unzip(source)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(4)
		}
	} else {
		documents = []io.ReadCloser{source}
	}
	for _, document := range documents {
		defer document.Close()
		if *parse {
			data, err := parser(document)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(5)
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "    ")
			err = encoder.Encode(data)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(6)
			}
		} else {
			buffer := make([]byte, 1024)
			for {
				n, err := document.Read(buffer)
				if n >= 0 {
					buffer = buffer[:n]
					_, e := os.Stdout.Write(buffer)
					if e != nil {
						fmt.Fprintln(os.Stderr, e.Error())
						os.Exit(8)
					}
				}
				if err == io.EOF {
					return
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					os.Exit(9)
				}
			}
		}
	}
}
