package bencode

import (
	"errors"
	"fmt"
	"github.com/firba1/downpour/bencode/bare"
	"io"
	"reflect"
	"strconv"
)

/*
Dict marshalls the bencoded data in r into the struct it.
Any extra values that don't have corresponding keys are returned in as a map[string]interface{}, each value contianing output from bare.Value()

Dict panics if it is not a pointer to struct.
Dict panics if the data in r is malformed.
*/
func Dict(it interface{}, r io.ByteScanner) map[string]interface{} {
	ptr := reflect.ValueOf(it)
	if ptr.Kind() != reflect.Ptr {
		panic("Must be a pointer (to a struct)!")
	}

	return dict(ptr.Elem(), r)
}

func dict(it reflect.Value, r io.ByteScanner) map[string]interface{} {
	extras := make(map[string]interface{})

	if it.Kind() != reflect.Struct {
		panic(fmt.Sprint("Bencoded dicts can only be put in go structs!"))
	}

	// check for the start of a dict: 'd'
	char, err := r.ReadByte()
	if err != nil {
		panic(err.Error())
	} else if char != 'd' {
		panic(fmt.Sprintf("Unexpected character '%v', expected 'd'", char))
	}

	for char, err = r.ReadByte(); char != 'e'; char, err = r.ReadByte() {
		if err != nil {
			panic(err)
		}
		// unread byte so it can be read by the appropriate parsing function
		r.UnreadByte()

		// read the dictionary key
		key := str(r)
		// find the appropriate struct field
		field, err := getByTag(it, key)
		if err != nil {
			extras[key] = bare.Value(r)
		}

		// try to marshal based on the field's Kind
		switch field.Kind() {
		case reflect.String:
			val := str(r)
			field.SetString(val)
		case reflect.Int:
			val := int64(integer(r))
			field.SetInt(val)
		case reflect.Struct:
			ex := dict(field, r)
			extras[key] = ex
		case reflect.Slice:
			field.Set(list(field, r))
		}
	}
	return extras
}

func getByTag(it reflect.Value, key string) (v reflect.Value, err error) {
	typ := it.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == key || field.Tag.Get("bencode") == key {
			v = it.Field(i)
			return
		}
	}
	err = errors.New(fmt.Sprintf("No field matches key \"%v\"", key))
	return
}

func List(l interface{}, r io.ByteScanner) {
	listVal := reflect.ValueOf(l)
	if listVal.Elem().Kind() != reflect.Slice {
		panic("Cannot marshal a list into something that's not a slice!")
	}
	listVal.Elem().Set(list(listVal.Elem(), r))
}

func list(l reflect.Value, r io.ByteScanner) reflect.Value {
	char, err := r.ReadByte()
	if err != nil {
		panic(err.Error())
	} else if char != 'l' {
		panic("Expected 'l'")
	}

	for char, err = r.ReadByte(); char != 'e'; char, err = r.ReadByte() {
		if err != nil {
			panic(err)
		}
		r.UnreadByte()
		switch l.Type().Elem().Kind() {
		case reflect.String:
			l = reflect.Append(l, reflect.ValueOf(str(r)))
		case reflect.Int:
			l = reflect.Append(l, reflect.ValueOf(integer(r)))
		case reflect.Struct:
			val := reflect.New(l.Type().Elem()).Elem()
			dict(val, r)
			l = reflect.Append(l, val)
		case reflect.Slice:
			v := list(reflect.Zero(l.Type().Elem()), r)
			l = reflect.Append(l, v)
		}
	}
	return l
}

func str(r io.ByteScanner) string {
	lenstr := []byte("")
	for c, err := r.ReadByte(); c != ':'; c, err = r.ReadByte() {
		if err != nil {
			panic(err.Error())
		}

		if !isDigit(c) {
			panic(fmt.Sprintf("'%c' should be a digit or ':'!", c))
		}

		lenstr = append(lenstr, c)
	}

	length, err := strconv.Atoi(string(lenstr))
	if err != nil {
		panic(err.Error())
	}

	bytes := []byte{}
	for i := 0; i < length; i++ {
		c, err := r.ReadByte()
		if err != nil {
			panic(err.Error())
		}

		bytes = append(bytes, c)
	}

	return string(bytes)
}

func integer(r io.ByteScanner) int {
	c, err := r.ReadByte()
	if err != nil {
		panic(err.Error())
	} else if c != 'i' {
		panic("Expected 'i'")
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

func isDigit(c byte) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}
