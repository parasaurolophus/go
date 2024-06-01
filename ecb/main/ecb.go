// Copyright Kirk Rader 2024

package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"parasaurolophus/go/ecb"
	"parasaurolophus/go/utilities"
)

func encode(data ecb.Data) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(data)
	if err != nil {
		panic(err.Error())
	}
}

func raw(reader io.Reader) {
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			buffer = buffer[:n]
			os.Stdout.Write(buffer)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}
	}
}

func csvHandler(parse bool) utilities.ZipHandler {
	return func(entry *zip.File) (err error) {
		readCloser, err := entry.Open()
		if err != nil {
			panic(err.Error())
		}
		defer readCloser.Close()
		if parse {
			data, err := ecb.ParseCSV(readCloser)
			if err != nil {
				panic(err.Error())
			}
			encode(data)
		} else {
			raw(readCloser)
		}
		return
	}
}

func xml(parse bool, reader io.Reader) {
	if parse {
		data, err := ecb.ParseXML(reader)
		if err != nil {
			panic(err.Error())
		}
		encode(data)
	} else {
		raw(reader)
	}
}

// Invoke ecb.Fetch manually, to support interactive debugging.
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprint(os.Stderr, r)
			fmt.Fprintln(os.Stderr)
			os.Exit(5)
		}
	}()
	format := flag.String("format", "csv", "csv or xml")
	version := flag.String("version", "daily", "daily, ninety or historical")
	parse := flag.Bool("parse", false, "true or false")
	flag.Parse()
	var url string
	if *format == "csv" {
		switch *version {
		case "daily":
			url = ecb.DAILY_CSV_URL
		case "historical":
			url = ecb.HISTORICAL_CSV_URL
		default:
			fmt.Fprintf(os.Stderr, `"%s" is not a valid version for csv`, *version)
			fmt.Fprintln(os.Stderr)
			os.Exit(1)
		}
	} else if *format == "xml" {
		switch *version {
		case "daily":
			url = ecb.DAILY_XML_URL
		case "historical":
			url = ecb.HISTORICAL_XML_URL
		case "ninety":
			url = ecb.NINETY_DAY_XML_URL
		default:
			fmt.Fprintf(os.Stderr, `"%s" is not a valid version for csv`, *version)
			fmt.Fprintln(os.Stderr)
			os.Exit(2)
		}
	} else {
		fmt.Fprintf(os.Stderr, `"%s" is not a valid format`, *format)
		fmt.Fprintln(os.Stderr)
		os.Exit(2)
	}
	source, err := utilities.Fetch(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(3)
	}
	defer source.Close()
	switch *format {
	case "csv":
		utilities.ForZipReader(csvHandler(*parse), source)
	case "xml":
		xml(*parse, source)
	default:
		fmt.Fprintf(os.Stderr, `unsupported format "%s"`, *format)
		fmt.Fprintln(os.Stderr)
		os.Exit(4)
	}
}
