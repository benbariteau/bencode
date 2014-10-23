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
	case reflect.Struct:
	}
	return []byte{}
}
