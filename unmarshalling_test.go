package bencode

import (
	"bytes"
	"reflect"
	"testing"
)

type testcase struct {
	in  string
	out int
}

func TestConsumeInt(t *testing.T) {
	tests := []testcase{
		testcase{"i2e", 2},
		testcase{"i42e", 42},
		testcase{"i666e", 666},
	}
	for _, test := range tests {
		out := consumeInt(bytes.NewBuffer([]byte(test.in)))
		if i := out.Interface(); !reflect.DeepEqual(i, test.out) {
			t.Error("Expecting", test.out, "got", i)
		}
	}
}
