package bencode

import (
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
	}

	for _, test := range tests {
		out := string(Marshal(test.in))
		if out != test.out {
			t.Error("Expecting", test.out, "got", out)
		}
	}
}
