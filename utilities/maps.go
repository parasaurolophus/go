// Copyright Kirk Rader 2024

package utilities

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/google/uuid"
)

// Return a shallow copy of m, with only the keys specified by keepKeys.
func FilterMapKeys[K comparable, V any](m map[K]V, keepKeys ...K) map[K]V {
	var result = map[K]V{}
	for k, v := range m {
		if slices.Contains(keepKeys, k) {
			result[k] = v
		}
	}
	return result
}

// Interpret the value of k in m as a float64, signalling an error if k is
// missing or of an unupported type.
func GetRequiredFloat(k string, m map[string]any) (float64, error) {
	return GetRequiredValue(k, m, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
}

// Interpret the value of k in m as a string, signalling an error if k is
// missing or of an unupported type.
func GetRequiredString[K comparable](k K, m map[K]any) (string, error) {
	// note that the parser function will not be called in this circumstance
	return GetRequiredValue[K, string](k, m, nil)
}

// Interpret the value of k in m as an uuid.UUID, signalling an error if k is
// missing or of an unupported type.
func GetRequiredUUID[K comparable](k K, m map[K]any) (uuid.UUID, error) {
	return GetRequiredValue(k, m, uuid.Parse)
}

// Interpret the value of k in m as a V, signalling an error if k is missing or
// of an unupported type.
func GetRequiredValue[K comparable, V any](k K, m map[K]any, parser func(string) (V, error)) (value V, err error) {
	v, found := m[k]
	if !found {
		err = fmt.Errorf("%v missing from %v", k, m)
		return
	}
	switch x := v.(type) {
	case V:
		value = x
	case string:
		value, err = parser(x)
	case fmt.Stringer:
		value, err = parser(x.String())
	default:
		err = fmt.Errorf("%v, of type %T, cannot be interpreted as a %T", x, x, value)
	}
	return
}

// Cast the value of k in m to V, signalling an error if the value is not of
// the specfied type.
func GetValue[K comparable, V any](k K, m map[K]any) (value V, found bool, err error) {
	v, found := m[k]
	if !found {
		return
	}
	value, ok := v.(V)
	if !ok {
		err = fmt.Errorf("expected %s to be of type %T, got %T", v, value, v)
	}
	return
}

// Return a map that combines the key / value pairs from all the given ones. The
// maps are processed in the given order. If the same key appears more than
// once, the value in the result will be the last one from the parameter list.
func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	var result = map[K]V{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
