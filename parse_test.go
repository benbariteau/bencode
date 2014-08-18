package bencode

import (
	"reflect"
	"strings"
	"testing"
)

type strtest struct {
	in  string
	out string
}

func TestStr(t *testing.T) {
	strs := []strtest{
		strtest{"3:str", "str"},
		strtest{"4:fart", "fart"},
		strtest{"11:fartbuttass", "fartbuttass"},
	}

	for _, test := range strs {
		out := str(strings.NewReader(test.in))
		if test.out != out {
			t.Error("Expected ", test.out, ", found ", out)
		}
	}
}

type inttest struct {
	in  string
	out int
}

func TestInteger(t *testing.T) {

	tests := []inttest{
		inttest{"i2e", 2},
		inttest{"i23e", 23},
	}

	for _, test := range tests {
		out := integer(strings.NewReader(test.in))
		if test.out != out {
			t.Error("Expected", test.out, " found", out)
		}
	}
}

func TestList(t *testing.T) {
	in1 := "l2:st3:stre"
	out1 := []string{"st", "str"}
	in2 := "li420ei666ee"
	out2 := []int{420, 666}
	in3 := "ld4:butti200eed4:butti404eee"
	out3 := []xinner{xinner{200}, xinner{404}}
	in4 := "ll4:fart4:buttel3:lol3:wutee"
	out4 := [][]string{[]string{"fart", "butt"}, []string{"lol", "wut"}}

	output := []string{}
	if List(&output, strings.NewReader(in1)); !reflect.DeepEqual(out1, output) {
		t.Error("expected", out1, "found", output)
	}

	output2 := []int{}
	if List(&output2, strings.NewReader(in2)); !reflect.DeepEqual(out2, output2) {
		t.Error("expected", out2, "found", output2)
	}

	output3 := []xinner{}
	if List(&output3, strings.NewReader(in3)); !reflect.DeepEqual(out3, output3) {
		t.Error("expected", out3, "found", output3)
	}

	output4 := [][]string{}
	if List(&output4, strings.NewReader(in4)); !reflect.DeepEqual(out4, output4) {
		t.Error("expected", out4, "found", output4)
	}
}

type strct struct {
	S string
	I int
	F inner
}

type inner struct {
	Fart string
	Shit xinner
}

type xinner struct {
	Butt int `bencode:"butt"`
}

type thing struct {
	Thingy []int
}

func TestDict(t *testing.T) {
	in := strct{}
	out := strct{"test", 420, inner{"yi", xinner{666}}}
	Dict(&in, strings.NewReader("d1:S4:test1:Ii420e1:Fd4:Fart2:yi4:Shitd4:butti666eeee"))
	if !reflect.DeepEqual(in, out) {
		t.Error("Expected", out, "found", in)
	}

	in2 := thing{}
	out2 := thing{[]int{123, 248}}
	r := strings.NewReader("d6:Thingyli123ei248eee")
	if Dict(&in2, r); !reflect.DeepEqual(in2, out2) {
		t.Error("Expected", out2, "found", in2)
	}
}
