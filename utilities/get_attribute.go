// Copyright 2024 Kirk Rader

package utilities

import (
	"fmt"
	"reflect"
	"strconv"
)

// Type constraint for arbitrary numeric conversions. Note that this excludes
// complex64 and complex128 because they do not support direct casting to and
// from the other numeric types.
type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// Convert the value of the specified key in the given map to the specified type.
func GetAttribute[Value any](m map[string]any, key string) (value Value, err error) {

	var (
		v  any
		ok bool
	)

	if v, ok = m[key]; !ok {
		err = fmt.Errorf("no value for %s in %v", key, m)
		return
	}

	if value, ok = v.(Value); !ok {
		err = fmt.Errorf("value of %s in %v is %v, not a %T", key, m, v, value)
	}

	return
}

// Convert the value of the specified key in the given map to the specified type of number.
func GetNumericAttribute[Value Number](m map[string]any, key string) (value Value, err error) {

	var (
		v  any
		ok bool
	)

	if v, ok = m[key]; !ok {
		err = fmt.Errorf("no value for %s in %v", key, m)
		return
	}

	switch vv := v.(type) {

	case int:
		value = Value(vv)

	case int8:
		value = Value(vv)

	case int16:
		value = Value(vv)

	case int32:
		value = Value(vv)

	case int64:
		value = Value(vv)

	case uint:
		value = Value(vv)

	case uint8:
		value = Value(vv)

	case uint16:
		value = Value(vv)

	case uint32:
		value = Value(vv)

	case uint64:
		value = Value(vv)

	case float32:
		value = Value(vv)

	case float64:
		value = Value(vv)

	case string:
		value, err = ParseNumber[Value](vv)

	case fmt.Stringer:
		value, err = ParseNumber[Value](vv.String())

	default:
		err = fmt.Errorf("value of %s in %v is of unsupported type %T", key, m, vv)
	}

	return
}

// Parse the given string as the specified type of number.
func ParseNumber[Value Number](s string) (value Value, err error) {

	t := reflect.TypeOf(value)

	switch t.Kind() {

	case reflect.Int8:
		var i int64
		if i, err = strconv.ParseInt(s, 10, 8); err == nil {
			value = Value(i)
		}

	case reflect.Int16:
		var i int64
		if i, err = strconv.ParseInt(s, 10, 16); err == nil {
			value = Value(i)
		}

	case reflect.Int32:
		var i int64
		if i, err = strconv.ParseInt(s, 10, 32); err == nil {
			value = Value(i)
		}

	case reflect.Int64:
		var i int64
		if i, err = strconv.ParseInt(s, 10, 64); err == nil {
			value = Value(i)
		}

	case reflect.Int:
		var i int64
		if i, err = strconv.ParseInt(s, 10, strconv.IntSize); err == nil {
			value = Value(i)
		}

	case reflect.Uint8:
		var i uint64
		if i, err = strconv.ParseUint(s, 10, 8); err == nil {
			value = Value(i)
		}

	case reflect.Uint16:
		var i uint64
		if i, err = strconv.ParseUint(s, 10, 16); err == nil {
			value = Value(i)
		}

	case reflect.Uint32:
		var i uint64
		if i, err = strconv.ParseUint(s, 10, 32); err == nil {
			value = Value(i)
		}

	case reflect.Uint64:
		var i uint64
		if i, err = strconv.ParseUint(s, 10, 64); err == nil {
			value = Value(i)
		}

	case reflect.Uint:
		var i uint64
		if i, err = strconv.ParseUint(s, 10, strconv.IntSize); err == nil {
			value = Value(i)
		}

	case reflect.Float32:
		var f float64
		if f, err = strconv.ParseFloat(s, 32); err == nil {
			value = Value(f)
		}

	default:
		var f float64
		if f, err = strconv.ParseFloat(s, 64); err == nil {
			value = Value(f)
		}
	}

	return
}
