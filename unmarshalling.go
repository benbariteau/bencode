package bencode

import (
	"bytes"
	"reflect"
)

func Unmarshal(data []byte, v interface{}) {
	container := reflect.ValueOf(v)
	buffer := bytes.NewBuffer(data)
}
