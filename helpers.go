package bencode

import (
	"fmt"
	"reflect"
)

type structHolder interface {
	getField(name string) reflect.Value
}

type realStructHolder struct {
	struc *reflect.Value
}

func (h realStructHolder) getField(name string) reflect.Value {
	structType := h.struc.Type()
	for i := 0; i < h.struc.NumField(); i++ {
		field := structType.Field(i)
		if field.PkgPath != "" { // if field not exported, it can't be set
			continue
		}
		tagName := field.Tag.Get("bencode")
		if tagName != "" && name == tagName {
			return h.struc.Field(i)
		} else if name == field.Name {
			return h.struc.Field(i)
		}
	}
	return reflect.Value{}
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
