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
	consumeValue(container.Elem(), buffer)
}

func consumeValue(variable reflect.Value, buffer *bytes.Buffer) {
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
		consumeInt(variable, buffer)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		consumeString(variable, buffer)
	case 'l':
		consumeList(variable, buffer)
	case 'd':
		consumeDict(variable, buffer)
	default:
		panic("Invalid thing")
	}
}

func consumeInt(variable reflect.Value, buffer *bytes.Buffer) {
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
	variable.SetInt(int64(value))
}

func consumeString(variable reflect.Value, buffer *bytes.Buffer) {
	lengthString, err := buffer.ReadString(':')
	if err != nil {
		panic("Unable to read string length: " + err.Error())
	}

	//remove trailing ':'
	lengthString = lengthString[:len(lengthString)-1]

	length, err := strconv.Atoi(lengthString)
	if err != nil {
		panic("Unable to convert number:" + err.Error())
	}

	bytes := buffer.Next(length)
	if len(bytes) < length {
		panic(fmt.Sprint("Expecting string of length", length, "got", len(bytes)))
	}
	variable.SetString(string(bytes))
}

func consumeList(variable reflect.Value, buffer *bytes.Buffer) {
	char, err := buffer.ReadByte()
	if err != nil {
		panic("Unable to read next byte:" + err.Error())
	}

	if char != 'l' {
		panic(fmt.Sprintf("Expecting 'l', found '%v'", char))
	}

	slice := variable
	for {
		char, err := buffer.ReadByte()
		if err != nil {
			panic("Unable to read next byte:" + err.Error())
		}

		if char == 'e' {
			break
		}

		err = buffer.UnreadByte()
		if err != nil {
			panic("Unable to read next byte:" + err.Error())
		}

		value := reflect.New(variable.Type().Elem()).Elem()
		consumeValue(value, buffer)
		slice = reflect.Append(slice, value)
	}
	variable.Set(slice)
}

func consumeDict(variable reflect.Value, buffer *bytes.Buffer) {
	char, err := buffer.ReadByte()
	if err != nil {
		panic("Unable to read next byte:" + err.Error())
	}

	if char != 'd' {
		panic(fmt.Sprintf("Expecting 'd', found '%v'", char))
	}

	for {
		char, err := buffer.ReadByte()
		if err != nil {
			panic("Unable to read next byte:" + err.Error())
		}

		if char == 'e' {
			break
		}

		err = buffer.UnreadByte()
		if err != nil {
			panic("Unable to read next byte:" + err.Error())
		}

		key := reflect.New(reflect.TypeOf("")).Elem()
		consumeString(key, buffer)
		field := variable.FieldByName(key.Interface().(string))
		consumeValue(field, buffer)
	}
}
