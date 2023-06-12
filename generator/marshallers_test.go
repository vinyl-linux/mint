package generator

import (
	"testing"

	"github.com/vinyl-linux/mint/parser"
)

func TestGenerator_marshallSliceArray(t *testing.T) {
	for _, test := range []struct {
		name   string
		ae     parser.AnnotatedEntry
		expect string
	}{
		{"Non-fixed length slice writes length", nonFixedLengthSlice, `func (sf TestType) marshallSomeStringSlice(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.SomeStringSlice))
	for i := range f {
		f[i] = mint.NewStringScalar(sf.SomeStringSlice[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}`},
		{"Fixed length slice writes length", fixedLengthSlice, `func (sf TestType) marshallSomeStringSlice(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.SomeStringSlice))
	for i := range f {
		f[i] = mint.NewStringScalar(sf.SomeStringSlice[i])
	}
	return mint.NewSliceCollection(f, true).Marshall(w)
}`},
		{"Bad input does nothing", parser.AnnotatedEntry{}, ""},
		{"Non-slice returns nothing", parser.AnnotatedEntry{Field: parser.Field{DataType: &parser.DataType{}}}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			g := new(Generator)
			received := codeToString(g.marshallSliceArray("TestType", test.ae))

			if test.expect != received {
				t.Errorf("expected\n%s\nreceived\n%s", test.expect, received)
			}
		})
	}
}

func TestGenerator_marshallMap(t *testing.T) {
	g := new(Generator)

	expect := `func (sf TestType) marshallSomeStringSlice(w io.Writer) (err error) {
	f := make(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer)
	for k, v := range sf.SomeStringSlice {
		f[mint.NewStringScalar(k)] = mint.NewInt64Scalar(v)
	}
	return mint.NewMapCollection(f).Marshall(w)
}`
	received := codeToString(g.marshallMap("TestType", mapEntry))

	if expect != received {
		t.Errorf("expected\n%s\nreceived\n%s", expect, received)
	}

}
