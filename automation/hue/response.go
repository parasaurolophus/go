// Copyright 2024 Kirk Rader

package hue

// HTTP response payload structure
type Response struct {
	Data   []Item `json:"data"`
	Errors []any  `json:"errors"`
}
