// Copyright Kirk Rader 2024

package utilities

import (
	"fmt"
	"reflect"
)

// Return a shallow copy of from, initializing only those fields named by
// keepFields. Note that the first parameter and return type is specified as
// 'any' due to limitations in Go's semantics.
func FilterStructFields(from any, keepFields ...string) (any, error) {
	typ := reflect.TypeOf(from)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported type %s", typ.String())
	}
	toPointer := reflect.New(typ)
	toValue := reflect.Indirect(toPointer)
	fromValue := reflect.ValueOf(from)
	for _, fieldName := range keepFields {
		toField := toValue.FieldByName(fieldName)
		if toField.CanSet() {
			toField.Set(fromValue.FieldByName(fieldName))
		}
	}
	return toValue.Interface(), nil
}
