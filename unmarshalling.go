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
		return consumeString(buffer)
	case 'l':
		return consumeList(buffer)
	case 'd':
		//TODO dict
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

func consumeString(buffer *bytes.Buffer) reflect.Value {
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
	return reflect.ValueOf(string(bytes))
}

func consumeList(buffer *bytes.Buffer) reflect.Value {
	char, err := buffer.ReadByte()
	if err != nil {
		panic("Unable to read next byte:" + err.Error())
	}

	if char != 'l' {
		panic(fmt.Sprintf("Expecting 'l', found '%v'", char))
	}

	//assumes list is homogenous
	firstValue := consumeValue(buffer)
	slice := reflect.Zero(reflect.SliceOf(firstValue.Type()))
	slice = reflect.Append(slice, firstValue)

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

		value := consumeValue(buffer)
		slice = reflect.Append(slice, value)
	}
	return slice
}
