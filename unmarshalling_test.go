package bencode

import (
	"bytes"
	"reflect"
	"testing"
)

type testcase struct {
	in  string
	out interface{}
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
	}

	for _, test := range tests {
		out := reflect.New(reflect.TypeOf(test.out)).Elem()
		parseDict(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}
