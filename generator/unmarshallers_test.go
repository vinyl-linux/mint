package generator

import (
	"testing"

	"github.com/vinyl-linux/mint/parser"
)

var (
	nonFixedLengthUnmarshaller = `func (sf *TestType) unmarshallSomeStringSlice(r io.Reader) (err error) {
	f := mint.NewSliceCollection(nil, false)
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	f.V = make([]mint.MarshallerUnmarshallerValuer, f.Len())
	for i := range f.V {
		f.V[i] = mint.NewStringScalar("")
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.SomeStringSlice = make([]string, f.Len())
	for i, v := range f.Value().([]mint.MarshallerUnmarshallerValuer) {
		sf.SomeStringSlice[i] = v.Value().(string)
	}
	return
}`

	fixedLengthUnmarshaller = `func (sf *TestType) unmarshallSomeStringSlice(r io.Reader) (err error) {
	f := mint.NewSliceCollection(make([]mint.MarshallerUnmarshallerValuer, 5), true)
	for i := range f.V {
		f.V[i] = mint.NewStringScalar("")
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.SomeStringSlice = [5]string{}
	for i, v := range f.Value().([]mint.MarshallerUnmarshallerValuer) {
		sf.SomeStringSlice[i] = v.Value().(string)
	}
	return
}`
)

func TestGenerator_unmarshallSliceArray(t *testing.T) {
	for _, test := range []struct {
		name   string
		ae     parser.AnnotatedEntry
		expect string
	}{
		{"Non-fixed length slice writes length", nonFixedLengthSlice, nonFixedLengthUnmarshaller},
		{"Fixed length slice writes length", fixedLengthSlice, fixedLengthUnmarshaller},
		{"Bad input does nothing", parser.AnnotatedEntry{}, ""},
		{"Non-slice returns nothing", parser.AnnotatedEntry{Field: parser.Field{DataType: &parser.DataType{}}}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			g := new(Generator)
			received := codeToString(g.unmarshallSliceArray("TestType", test.ae))

			if test.expect != received {
				t.Errorf("expected\n%s\nreceived\n%s", test.expect, received)
			}
		})

	}
}

func TestGenerator_unmarshallMap(t *testing.T) {
	g := new(Generator)

	expect := `func (sf *TestType) unmarshallSomeStringSlice(r io.Reader) (err error) {
	f := mint.NewMapCollection(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer{})
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	for i := 0; i < f.Len(); i++ {
		f.V[mint.NewStringScalar("")] = mint.NewInt64Scalar(int64(0))
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.SomeStringSlice = make(map[string]int64)
	for k, v := range f.Value().(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer) {
		sf.SomeStringSlice[k.Value().(string)] = v.Value().(int64)
	}
	return
}`
	received := codeToString(g.unmarshallMap("TestType", mapEntry))

	if expect != received {
		t.Errorf("expected\n%s\nreceived\n%s", expect, received)
	}

}

func TestGenerator_unmarshallScalar(t *testing.T) {
	g := new(Generator)

	expect := `func (sf *TestType) unmarshallATypeOfSomeType(r io.Reader) (err error) {
	f := mint.NewUuidScalar(v5.UUID{})
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.ATypeOfSomeType = f.Value().(v5.UUID)
	return
}`
	received := codeToString(g.unmarshallScalar("TestType", scalarEntry))

	if expect != received {
		t.Errorf("expected\n%s\nreceived\n%s", expect, received)
	}
}
