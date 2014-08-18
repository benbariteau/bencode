package bare

import (
	"reflect"
	"strings"
	"testing"
)

type dinout struct {
	in  string
	out map[string]interface{}
}

func TestDict(t *testing.T) {
	tests := []dinout{
		dinout{
			"d3:foo3:bar3:bazi666e4:fartli1ei2ei3ee4:buttd4:deep3:manee",
			map[string]interface{}{
				"foo":  "bar",
				"baz":  666,
				"fart": []interface{}{1, 2, 3},
				"butt": map[string]interface{}{
					"deep": "man",
				},
			},
		},
	}

	for _, test := range tests {
		if out := Dict(strings.NewReader(test.in)); !reflect.DeepEqual(test.out, out) {
			t.Errorf("Expecting '%v', got '%v'", test.out, out)
		}
	}
}

type linout struct {
	in  string
	out []interface{}
}

func TestList(t *testing.T) {
	tests := []linout{
		linout{"l4:fart4:butte", []interface{}{"fart", "butt"}},
		linout{"li420ei666ee", []interface{}{420, 666}},
		linout{"ll3:fooel3:baree", []interface{}{
			[]interface{}{"foo"},
			[]interface{}{"bar"},
		}},
		linout{
			"ld3:foo3:baree",
			[]interface{}{
				map[string]interface{}{"foo": "bar"},
			},
		},
	}

	for _, test := range tests {
		if out := List(strings.NewReader(test.in)); !reflect.DeepEqual(test.out, out) {
			t.Errorf("Expecting '%v', got '%v'", test.out, out)
		}
	}
}

type sinout struct {
	in  string
	out string
}

func TestString(t *testing.T) {
	tests := []sinout{
		sinout{"3:wut", "wut"},
		sinout{"4:fart", "fart"},
	}

	for _, test := range tests {
		if out := String(strings.NewReader(test.in)); test.out != out {
			t.Errorf("Expecting '%v', got '%v'", test.out, out)
		}
	}
}

type iinout struct {
	in  string
	out int
}

func TestInt(t *testing.T) {
	tests := []iinout{
		iinout{"i666e", 666},
	}

	for _, test := range tests {
		if out := Int(strings.NewReader(test.in)); test.out != out {
			t.Errorf("Expecting '%v', got '%v'", test.out, out)
		}
	}
}
