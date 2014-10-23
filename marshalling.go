package bencode

import (
	"fmt"
	"reflect"
	"strconv"
)

func Marshal(v interface{}) []byte {
	value := reflect.ValueOf(v)
	return convertValue(value)
}

func convertValue(value reflect.Value) []byte {
	switch value.Type().Kind() {
	case reflect.Int:
		return []byte("i" + strconv.Itoa(value.Interface().(int)) + "e")
	case reflect.String:
		stringValue := value.Interface().(string)
		return []byte(fmt.Sprintf("%v:%v", len([]byte(stringValue)), stringValue))
	case reflect.Slice:
		return convertSlice(value)
	case reflect.Struct:
	}
	return []byte{}
}

func convertSlice(value reflect.Value) (representation []byte) {
	representation = append(representation, 'l')
	for i := 0; i < value.Len(); i++ {
		valueRepresentation := convertValue(value.Index(i))
		representation = append(representation, valueRepresentation...)
	}
	representation = append(representation, 'e')
	return
}
