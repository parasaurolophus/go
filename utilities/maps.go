// Copyright Kirk Rader 2024

package utilities

import (
	"slices"
)

func FilterMapKeys[K comparable, V any](m map[K]V, keepKeys ...K) map[K]V {

	var result = map[K]V{}

	for k, v := range m {

		if slices.Contains(keepKeys, k) {

			result[k] = v
		}
	}

	return result
}

func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {

	var result = map[K]V{}

	for _, m := range maps {

		for k, v := range m {

			result[k] = v
		}
	}

	return result
}
