package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Unmarshal(data []byte, v interface{}) (err error) {

	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case string:
				err = errors.New(r)
			case error:
				err = r
			}
		}
	}()

	container := reflect.ValueOf(v)
	buffer := bytes.NewBuffer(data)
	consumeValue(container.Elem(), buffer)
	return nil
}

var parseMapInitialized bool = false
var regularParseMap map[byte]parserFunc
var ignoreParseMap map[byte]parserFunc

type parserFunc func(reflect.Value, *bytes.Buffer)

func generateParserMap(intParser, stringParser, listParser, dictParser parserFunc) map[byte]parserFunc {
	parserMap := make(map[byte]parserFunc)
	parserMap['i'] = intParser
	parserMap['l'] = listParser
	parserMap['d'] = dictParser
	for i := 0; i <= 9; i++ {
		key := strconv.Itoa(i)[0]
		parserMap[key] = stringParser
	}
	return parserMap
}

func consumeValue(variable reflect.Value, buffer *bytes.Buffer) {

	if !parseMapInitialized {
		parseMapInitialized = true
		regularParseMap = generateParserMap(parseInt, parseString, parseList, parseDict)
		ignoreParseMap = generateParserMap(ignoreInt, ignoreString, ignoreList, ignoreDict)
	}

	parseMap := regularParseMap
	if !variable.IsValid() {
		parseMap = ignoreParseMap
	}

	char, err := buffer.ReadByte()
	if err != nil {
		//TODO replace with error type
		panic("Unable to read next byte:" + err.Error())
	}
	buffer.UnreadByte()

	handler, ok := parseMap[char]
	if !ok {
		panic(fmt.Sprintf("Expecting 'i', 'l', 'd', or a digit (0-9), found, '%v'", char))
	}
	handler(variable, buffer)
}

func parseInt(variable reflect.Value, buffer *bytes.Buffer) {
	variable.SetInt(int64(consumeInt(buffer)))
}

func ignoreInt(variable reflect.Value, buffer *bytes.Buffer) {
	consumeInt(buffer)
}

func consumeInt(buffer *bytes.Buffer) int {
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
	return value
}

func parseString(variable reflect.Value, buffer *bytes.Buffer) {
	variable.SetString(consumeString(buffer))
}

func ignoreString(variable reflect.Value, buffer *bytes.Buffer) {
	consumeString(buffer)
}

func consumeString(buffer *bytes.Buffer) string {
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
	return string(bytes)
}

func parseList(variable reflect.Value, buffer *bytes.Buffer) {
	listBuffer := realBuffer{&variable}
	parseListHelper(listBuffer, variable, buffer)
}

func parseListHelper(listBuffer sliceBuffer, variable reflect.Value, buffer *bytes.Buffer) {
	consumeList(listBuffer, buffer)
	variable.Set(listBuffer.value())
}

func ignoreList(variable reflect.Value, buffer *bytes.Buffer) {
	consumeList(fakeBuffer{}, buffer)
}

func consumeList(listBuffer sliceBuffer, buffer *bytes.Buffer) {
	char, err := buffer.ReadByte()
	if err != nil {
		panic("Unable to read next byte:" + err.Error())
	}

	if char != 'l' {
		panic(fmt.Sprintf("Expecting 'l', found '%v'", char))
	}

	for {
		char, err := buffer.ReadByte()
		if err != nil {
			panic("Unable to read next byte:" + err.Error())
		}

		if char == 'e' {
			break
		}

		buffer.UnreadByte()

		value := listBuffer.newValue()
		consumeValue(value, buffer)
		listBuffer.push(value)
	}
}

func parseDict(variable reflect.Value, buffer *bytes.Buffer) {
	thing := realStructHolder{&variable}
	consumeDict(thing, buffer)
}

func ignoreDict(variable reflect.Value, buffer *bytes.Buffer) {
	consumeDict(fakeStructHolder{}, buffer)
}

func consumeDict(thing structHolder, buffer *bytes.Buffer) {
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

		buffer.UnreadByte()

		key := reflect.New(reflect.TypeOf("")).Elem()
		parseString(key, buffer)
		field := thing.getField(key.Interface().(string))
		consumeValue(field, buffer)
	}
}
