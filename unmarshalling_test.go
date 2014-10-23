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

func TestConsumeInt(t *testing.T) {
	tests := []testcase{
		testcase{"i2e", 2},
		testcase{"i42e", 42},
		testcase{"i666e", 666},
	}
	for _, test := range tests {
		var o int
		out := reflect.ValueOf(&o).Elem()
		consumeInt(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}

func TestConsumeString(t *testing.T) {
	tests := []testcase{
		testcase{"1:s", "s"},
		testcase{"4:butt", "butt"},
		testcase{"11:buttfartass", "buttfartass"},
	}
	for _, test := range tests {
		var o string
		out := reflect.ValueOf(&o).Elem()
		consumeString(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}

func TestConsumeList(t *testing.T) {
	tests := []testcase{
		testcase{"li2ei42ei666ee", []int{2, 42, 666}},
		testcase{"l1:a1:b1:ce", []string{"a", "b", "c"}},
		testcase{"lli2eeli42eee", [][]int{[]int{2}, []int{42}}},
		testcase{"ld1:A4:butted1:A4:fartee", []struct{ A string }{
			struct{ A string }{"butt"},
			struct{ A string }{"fart"},
		}},
	}

	for _, test := range tests {
		out := reflect.New(reflect.TypeOf(test.out)).Elem()
		consumeList(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}

func TestConsumeDict(t *testing.T) {
	tests := []testcase{
		testcase{"d1:Ai42e1:B3:xyze", struct {
			A int
			B string
		}{A: 42, B: "xyz"}},
	}

	for _, test := range tests {
		out := reflect.New(reflect.TypeOf(test.out)).Elem()
		consumeDict(out, bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}
