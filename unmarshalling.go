package bencode

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

func Unmarshal(data []byte, v interface{}) {
	container := reflect.ValueOf(v)
	buffer := bytes.NewBuffer(data)
	value := consumeValue(buffer)
	container.Set(value)
}

func consumeValue(buffer *bytes.Buffer) reflect.Value {
	char, err := buffer.ReadByte()
	if err != nil {
		//TODO replace with error type
		panic("Unable to read next byte:" + err.Error())
	}
	err = buffer.UnreadByte()
	if err != nil {
		panic("Unable to read next byte:" + err.Error())
	}

	switch char {
	case 'i':
		return consumeInt(buffer)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		//TODO string
	case 'd':
		//TODO dict
	case 'l':
		//TODO list
	default:
		panic("Invalid thing")
	}
	return reflect.ValueOf(nil)
}

func consumeInt(buffer *bytes.Buffer) reflect.Value {
	char, err := buffer.ReadByte()
	if err != nil {
		panic("Unable to read next byte:" + err.Error())
	}

	if char != 'i' {
		panic(fmt.Sprintf("Expecting 'i', found '%v'", char))
	}

	decimalString, err := buffer.ReadString('e')
	if err != nil {
		panic("Unable to read next byte:" + err.Error())
	}

	//remove trailing 'e'
	decimalString = decimalString[:len(decimalString)-1]

	value, err := strconv.Atoi(decimalString)
	if err != nil {
		panic("Unable to convert number:" + err.Error())
	}
	return reflect.ValueOf(value)
}
