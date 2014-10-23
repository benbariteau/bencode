package bencode

import (
	"reflect"
	"testing"
)

type mtestcase struct {
	in  interface{}
	out string
}

func TestMarshal(t *testing.T) {
	tests := []mtestcase{
		mtestcase{2, "i2e"},
		mtestcase{"fart", "4:fart"},
		mtestcase{[]string{"fart", "butt"}, "l4:fart4:butte"},
	}

	for _, test := range tests {
		out := string(Marshal(test.in))
		if out != test.out {
			t.Error("Expecting", test.out, "got", out)
		}
	}
}

func TestConvertSlice(t *testing.T) {
	type twofield struct {
		A int
		B string
	}
	tests := []mtestcase{
		mtestcase{[]int{2, 42, 666}, "li2ei42ei666ee"},
		mtestcase{[]string{"fart", "butt"}, "l4:fart4:butte"},
		mtestcase{[][]int{[]int{2}, []int{42}}, "lli2eeli42eee"},
		mtestcase{
			[]twofield{twofield{420, "blaze it"}, twofield{360, "no scope"}},
			"ld1:Ai420e1:B8:blaze ited1:Ai360e1:B8:no scopeee",
		},
	}
	for _, test := range tests {
		out := string(convertSlice(reflect.ValueOf(test.in)))
		if out != test.out {
			t.Error("Expecting", test.out, "got", out)
		}
	}
}

func TestConvertDict(t *testing.T) {
	tests := []mtestcase{
		mtestcase{struct{ A string }{A: "butt"}, "d1:A4:butte"},
		mtestcase{struct{ A []int }{A: []int{42, 666}}, "d1:Ali42ei666eee"},
		mtestcase{struct{ A struct{ A string } }{A: struct{ A string }{A: "butt"}}, "d1:Ad1:A4:buttee"},
	}
	for _, test := range tests {
		out := string(convertDict(reflect.ValueOf(test.in)))
		if out != test.out {
			t.Error("Expecting", test.out, "got", out)
		}
	}
}
