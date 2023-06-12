package generator

import (
	"testing"
)

func TestToGoFunc(t *testing.T) {
	for _, test := range []struct {
		name     string
		custom   bool
		function string
		expect   string
	}{
		{"Non-custom function is prefixed with mint, camel cased", false, "do_something", "mint.DoSomething"},
		{"Custom function is prefixed with sf, camel cased", true, "a_b_c_de", "sf.ABCDe"},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := toGoFuncName(test.custom, test.function).GoString()
			if test.expect != received {
				t.Errorf("expected %q, received %q", test.expect, received)
			}
		})
	}
}

func TestMarshallerFuncName(t *testing.T) {
	expect := "marshallField"
	received := marshallerFuncName("Field")

	if expect != received {
		t.Errorf("expected %q, received %q", expect, received)
	}
}

func TestUnmarshallerFuncName(t *testing.T) {
	expect := "unmarshallField"
	received := unmarshallerFuncName("Field")

	if expect != received {
		t.Errorf("expected %q, received %q", expect, received)
	}
}
