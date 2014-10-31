package bencode

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

type testcase struct {
	in  string
	out interface{}
}

func TestUnmarshal(t *testing.T) {
	tests := []testcase{
		testcase{"i2e", 2},
		testcase{"1:a", "a"},
		testcase{"li2ei42ei666ee", []int{2, 42, 666}},
		testcase{"d1:Ai42e1:B3:xyze", struct {
			A int
			B string
		}{A: 42, B: "xyz"}},
	}

	for _, test := range tests {
		out := reflect.New(reflect.TypeOf(test.out))
		err := Unmarshal([]byte(test.in), out.Interface())
		if err != nil {
			t.Error("Error while running test:", err)
		} else if i := out.Elem().Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}

func TestInvalidUnmarshal(t *testing.T) {
	tests := []testcase{
		testcase{"a", 0},
		testcase{"", 0},
	}

	for _, test := range tests {
		out := reflect.New(reflect.TypeOf(test.out))
		err := Unmarshal([]byte(test.in), out.Interface())
		if err == nil {
			t.Errorf("Expecting error for input: \"%v\"", test.in)
		}
	}
}

func TestNonpointerUnmarshal(t *testing.T) {
	err := Unmarshal([]byte{}, 1)
	if err == nil {
		t.Errorf("Expecting error on non-pointer input to Unmarshal")
	}
}

func TestConsumeValue(t *testing.T) {
	tests := []testcase{
		testcase{"i2e", 2},
		testcase{"1:a", "a"},
		testcase{"li2ei42ei666ee", []int{2, 42, 666}},
		testcase{"d1:Ai42e1:B3:xyze", struct {
			A int
			B string
		}{A: 42, B: "xyz"}},
	}
	for _, test := range tests {
		out := reflect.New(reflect.TypeOf(test.out)).Elem()
		consumeValue(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}

func TestParseInt(t *testing.T) {
	tests := []testcase{
		testcase{"i2e", 2},
		testcase{"i42e", 42},
		testcase{"i666e", 666},
	}
	for _, test := range tests {
		var o int
		out := reflect.ValueOf(&o).Elem()
		parseInt(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}

func TestInvalidParseInt(t *testing.T) {
	tests := []testcase{
		testcase{"2e", 2},
		testcase{"i42", 42},
		testcase{"666", 666},
		testcase{"i2/e", 2},
		testcase{"", 0},
	}
	for _, test := range tests {
		var o int
		out := reflect.ValueOf(&o).Elem()
		err := errorCatcher(func() {
			parseInt(out, bytes.NewBuffer([]byte(test.in)))
		})
		if err == nil {
			t.Errorf("Expecting error for input: \"%v\"", test.in)
		}
	}
}

func TestParseString(t *testing.T) {
	tests := []testcase{
		testcase{"1:s", "s"},
		testcase{"4:butt", "butt"},
		testcase{"11:buttfartass", "buttfartass"},
	}
	for _, test := range tests {
		var o string
		out := reflect.ValueOf(&o).Elem()
		parseString(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}
func TestInvalidParseString(t *testing.T) {
	tests := []testcase{
		testcase{"4", ""},
		testcase{"4/:", ""},
		testcase{"5:fart", "fart"},
	}
	for _, test := range tests {
		var o string
		out := reflect.ValueOf(&o).Elem()
		err := errorCatcher(func() {
			parseString(out, bytes.NewBuffer([]byte(test.in)))
		})
		if err == nil {
			t.Errorf("Expecting error for input: \"%v\"", test.in)
		}
	}
}

func TestParseList(t *testing.T) {
	type onefield struct {
		A string
	}
	tests := []testcase{
		testcase{"li2ei42ei666ee", []int{2, 42, 666}},
		testcase{"l1:a1:b1:ce", []string{"a", "b", "c"}},
		testcase{"lli2eeli42eee", [][]int{[]int{2}, []int{42}}},
		testcase{"ld1:A4:butted1:A4:fartee", []onefield{
			onefield{"butt"},
			onefield{"fart"},
		}},
	}

	for _, test := range tests {
		out := reflect.New(reflect.TypeOf(test.out)).Elem()
		parseList(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}

func TestInvalidParseList(t *testing.T) {
	tests := []testcase{
		testcase{"", []int{}},
		testcase{"f", []int{}},
		testcase{"l", []int{}},
	}
	for _, test := range tests {
		out := reflect.ValueOf(reflect.TypeOf(test.out)).Elem()
		err := errorCatcher(func() {
			parseList(out, bytes.NewBuffer([]byte(test.in)))
		})
		if err == nil {
			t.Errorf("Expecting error for input: \"%v\"", test.in)
		}
	}
}

func TestParseDict(t *testing.T) {
	type twofield struct {
		A int
		B string
	}
	type liststruct struct {
		A []int
	}
	type onefield struct {
		B string
	}
	type structstruct struct {
		A onefield
	}
	type tagged struct {
		A string `bencode:"tag"`
	}
	type unexported struct {
		A int
		b int
	}
	tests := []testcase{
		testcase{"d1:Ai42e1:B3:xyze", twofield{
			A: 42, B: "xyz",
		}},
		testcase{"d1:Ali2ei42ei666eee", liststruct{
			A: []int{2, 42, 666},
		}},
		testcase{"d1:Ad1:B3:xyzee", structstruct{
			A: onefield{
				B: "xyz",
			},
		}},
		testcase{"d3:tag3:xyze", tagged{A: "xyz"}},
		testcase{"d3:tag3:xyz1:a1:be", tagged{A: "xyz"}},
		testcase{"d1:Ai42e1:bi666ee", unexported{A: 42}},
		testcase{"d1:B4:fart1:cd1:di666eee", onefield{"fart"}},
		testcase{"d1:B4:fart1:cli1ei2ei3eee", onefield{"fart"}},
	}

	for _, test := range tests {
		out := reflect.New(reflect.TypeOf(test.out)).Elem()
		parseDict(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}

func TestInvalidParseDict(t *testing.T) {
	tests := []testcase{
		testcase{"", []int{}},
		testcase{"f", []int{}},
		testcase{"d", []int{}},
	}
	for _, test := range tests {
		out := reflect.ValueOf(reflect.TypeOf(test.out)).Elem()
		err := errorCatcher(func() {
			parseDict(out, bytes.NewBuffer([]byte(test.in)))
		})
		if err == nil {
			t.Errorf("Expecting error for input: \"%v\"", test.in)
		}
	}
}

func errorCatcher(f func()) (err error) {
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
	f()
	return nil
}
