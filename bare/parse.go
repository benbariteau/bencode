package bare

import (
	"fmt"
	"io"
	"strconv"
)

func Value(r io.ByteScanner) interface{} {
	char, err := r.ReadByte()
	if err != nil {
		panic(err.Error())
	}
	r.UnreadByte()

	switch char {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return String(r)
	case 'i':
		return Int(r)
	case 'l':
		return List(r)
	case 'd':
		return Dict(r)
	}
	panic(fmt.Sprintf("Unexpected char '%c', not a valid bencoded value beginning", char))
}

func Dict(r io.ByteScanner) map[string]interface{} {
	c, err := r.ReadByte()
	if err != nil {
		panic(err.Error())
	} else if c != 'd' {
		panic(fmt.Sprintf("Expected 'd', got '%c'", c))
	}

	dict := make(map[string]interface{})

	for char, err := r.ReadByte(); char != 'e'; char, err = r.ReadByte() {
		if err != nil {
			panic(err.Error())
		}
		r.UnreadByte()

		key := String(r)
		val := Value(r)
		dict[key] = val
	}

	return dict
}

func List(r io.ByteScanner) (list []interface{}) {
	c, err := r.ReadByte()
	if err != nil {
		panic(err.Error())
	} else if c != 'l' {
		panic(fmt.Sprintf("Expected 'l', got '%c'", c))
	}

	for char, err := r.ReadByte(); char != 'e'; char, err = r.ReadByte() {
		if err != nil {
			panic(err.Error())
		}
		r.UnreadByte()
		list = append(list, Value(r))
	}
	return
}

func String(r io.ByteScanner) string {
	var lenbytestr []byte
	for char, err := r.ReadByte(); char != ':'; char, err = r.ReadByte() {
		if err != nil {
			panic(err.Error())
		}

		if !isDigit(char) {
			panic(fmt.Sprintf("Expecting digit, got '%v'", char))
		}
		lenbytestr = append(lenbytestr, char)
	}

	length, err := strconv.ParseInt(string(lenbytestr), 10, 0)
	if err != nil {
		panic(err.Error())
	}

	var strbytes []byte
	for i := 0; int64(i) < length; i++ {
		char, err := r.ReadByte()
		if err != nil {
			panic(err.Error())
		}

		strbytes = append(strbytes, char)
	}
	return string(strbytes)
}

func isDigit(c byte) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

func Int(r io.ByteReader) int {
	c, err := r.ReadByte()
	if err != nil {
		panic(err.Error())
	} else if c != 'i' {
		panic(fmt.Sprintf("Expected 'i', got '%c'", c))
	}

	numstr := []byte("")
	for c, err := r.ReadByte(); c != 'e'; c, err = r.ReadByte() {
		if err != nil {
			panic(err.Error())
		}

		if !isDigit(c) {
			panic(fmt.Sprintf("'%c' should be a digit or 'e'!", c))
		}

		numstr = append(numstr, c)
	}
	num, err := strconv.Atoi(string(numstr))
	if err != nil {
		panic(err.Error())
	}

	return num
}
