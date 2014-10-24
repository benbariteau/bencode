package bencode

import (
	"fmt"
	"reflect"
)

func getStructField(name string, variable reflect.Value) (reflect.Value, error) {
	structType := variable.Type()
	for i := 0; i < variable.NumField(); i++ {
		field := structType.Field(i)
		if field.PkgPath != "" { // if field not exported, it can't be set
			continue
		}
		tagName := field.Tag.Get("bencode")
		if tagName != "" && name == tagName {
			return variable.Field(i), nil
		} else if name == field.Name {
			return variable.Field(i), nil
		}
	}
	return reflect.Value{}, noFieldError{name, structType}
}

type noFieldError struct {
	fieldName  string
	structType reflect.Type
}

func (e noFieldError) Error() string {
	return fmt.Sprintf("Field '%v' no in struct type '%T'", e.fieldName, e.structType)
}

type sliceBuffer interface {
	newValue() reflect.Value
	push(value reflect.Value)
	value() reflect.Value
}

type realBuffer struct {
	slice *reflect.Value
}

func (s realBuffer) newValue() reflect.Value {
	return reflect.New(s.slice.Type().Elem()).Elem()
}

func (s realBuffer) push(value reflect.Value) {
	*s.slice = reflect.Append(*s.slice, value)
}

func (s realBuffer) String() string {
	return fmt.Sprint(s.slice.Interface())
}

func (s realBuffer) value() reflect.Value {
	return *s.slice
}
