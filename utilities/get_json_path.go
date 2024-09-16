// Copyright 2024 Kirk Rader

package utilities

import (
	"fmt"
)

// Return the value specified by the given path in the given map. For example,
// {"foo"} is equivalent to m["foo"] while {"foo", ""bar"} is equivalent to
// m["foo"]["bar"]. For paths with more than one entry, each intermediate
// container is assumed to be a map[string]any.
func GetJSONPath[Value any](

	m map[string]any,
	path ...string,

) (

	value Value,
	err error,

) {

	var (
		a  any
		ok bool
	)

	n := len(path)

	if n < 1 {
		err = fmt.Errorf("empty path")
		return
	}

	if a, ok = m[path[0]]; !ok {
		err = fmt.Errorf(`no value found for "%s" in %v`, path[0], m)
		return
	}

	if n > 1 {

		var c map[string]any
		if c, ok = a.(map[string]any); !ok {
			err = fmt.Errorf(`"%s" in %v is not a map`, path[0], m)
			return
		}
		path = path[1:]
		value, err = GetJSONPath[Value](c, path...)

	} else if value, ok = a.(Value); !ok {

		err = fmt.Errorf(
			`value %v, of type %T, found for "%s" in %v; %T expected`,
			a,
			a,
			path[0],
			m,
			value,
		)
	}

	return

}
