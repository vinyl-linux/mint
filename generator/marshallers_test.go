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
		{"Complex types are not passed to an initialiser", userDefinedSliceEntry, `func (sf TestType) marshallThingy(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.Thingy))
	for i := range f {
		f[i] = &(sf.Thingy[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}`},
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

	for _, test := range []struct {
		name   string
		ae     parser.AnnotatedEntry
		expect string
	}{
		{"map of builtin to builtin", mapEntry, `func (sf TestType) marshallSomeStringSlice(w io.Writer) (err error) {
	f := make(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer)
	for k, v := range sf.SomeStringSlice {
		f[mint.NewStringScalar(k)] = mint.NewInt64Scalar(v)
	}
	return mint.NewMapCollection(f).Marshall(w)
}`},
		{"map of builtin to complex", builtinToComplexMap, `func (sf TestType) marshallThingy(w io.Writer) (err error) {
	f := make(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer)
	for k, v := range sf.Thingy {
		f[mint.NewStringScalar(k)] = &(v)
	}
	return mint.NewMapCollection(f).Marshall(w)
}`},
		{"map of complex to builtin", complexToBuiltinMap, `func (sf TestType) marshallThingy(w io.Writer) (err error) {
	f := make(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer)
	for k, v := range sf.Thingy {
		f[&(k)] = mint.NewBoolScalar(v)
	}
	return mint.NewMapCollection(f).Marshall(w)
}`},
		{"map of complex to complex", complexToComplexMap, `func (sf TestType) marshallThingy(w io.Writer) (err error) {
	f := make(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer)
	for k, v := range sf.Thingy {
		f[&(k)] = &(v)
	}
	return mint.NewMapCollection(f).Marshall(w)
}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			g := new(Generator)
			received := codeToString(g.marshallMap("TestType", test.ae))

			if test.expect != received {
				t.Errorf("expected\n%s\nreceived\n%s", test.expect, received)
			}
		})
	}
}
