// Copyright Kirk Rader 2024

package utilities

import (
	"reflect"
)

func FilterStructFields(from any, keepFields ...string) any {

	typ := reflect.TypeOf(from)
	toPointer := reflect.New(typ)
	toValue := reflect.Indirect(toPointer)

	if typ.Kind() != reflect.Struct {
		return nil
	}

	fromValue := reflect.ValueOf(from)

	for _, fieldName := range keepFields {

		toField := toValue.FieldByName(fieldName)

		if toField.CanSet() {

			toField.Set(fromValue.FieldByName(fieldName))
		}
	}

	return toValue.Interface()
}
