// Copyright Kirk Rader 2024

package utilities

import (
	"slices"
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
