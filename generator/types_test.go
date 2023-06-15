package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/vinyl-linux/mint/parser"
)

var (
	nonFixedLengthSlice = parser.AnnotatedEntry{
		Field: parser.Field{
			Name: "SomeStringSlice",
			DataType: &parser.DataType{
				Slice: &parser.SliceType{
					Type: "string",
				},
			},
		},
	}

	fixedLengthSlice = parser.AnnotatedEntry{
		Transformations: []parser.Transformation{
			{
				IsCustom: true,
				Function: "floop",
			},
		},
		Field: parser.Field{
			Name: "SomeStringSlice",
			DataType: &parser.DataType{
				FixedSizeSlice: &parser.FixedSizedSliceType{
					Type: "string",
					Size: 5,
				},
			},
		},
	}

	mapEntry = parser.AnnotatedEntry{
		DocString: "An int64 for some reason",
		Transformations: []parser.Transformation{
			{
				IsCustom: true,
				Function: "treble_value",
			},
			{
				IsCustom: false,
				Function: "flipbits",
			},
		},
		Field: parser.Field{
			Name: "SomeStringSlice",
			DataType: &parser.DataType{
				Map: &parser.MapType{
					Key:   "string",
					Value: "int64",
				},
			},
		},
	}

	scalarEntry = parser.AnnotatedEntry{
		DocString: "ATypeOfSomeType is a uuid",
		Validations: []parser.Validation{
			{
				IsCustom: false,
				Function: "not_empty",
			},
			{
				IsCustom: true,
				Function: "blah_blah_blah-how-does-this.evaluate",
			},
		},
		Field: parser.Field{
			Name: "ATypeOfSomeType",
			DataType: &parser.DataType{
				Scalar: &parser.Scalar{
					Type: "uuid",
				},
			},
		},
	}

	userDefinedScalarEntry = parser.AnnotatedEntry{
		Field: parser.Field{
			Name: "Thingy",
			DataType: &parser.DataType{
				Scalar: &parser.Scalar{
					Type: "BlahType",
				},
			},
		},
	}

	userDefinedSliceEntry = parser.AnnotatedEntry{
		Field: parser.Field{
			Name: "Thingy",
			DataType: &parser.DataType{
				Slice: &parser.SliceType{
					Type: "BlahType",
				},
			},
		},
	}

	builtinToComplexMap = parser.AnnotatedEntry{
		Field: parser.Field{
			Name: "Thingy",
			DataType: &parser.DataType{
				Map: &parser.MapType{
					Key:   "string",
					Value: "BlahType",
				},
			},
		},
	}

	complexToBuiltinMap = parser.AnnotatedEntry{
		Field: parser.Field{
			Name: "Thingy",
			DataType: &parser.DataType{
				Map: &parser.MapType{
					Key:   "BlahType",
					Value: "bool",
				},
			},
		},
	}

	complexToComplexMap = parser.AnnotatedEntry{
		Field: parser.Field{
			Name: "Thingy",
			DataType: &parser.DataType{
				Map: &parser.MapType{
					Key:   "BlahType",
					Value: "BlahType",
				},
			},
		},
	}

	simpleType = parser.AnnotatedType{
		Name: "SomeTestType",
		Entries: []parser.AnnotatedEntry{
			nonFixedLengthSlice,
			fixedLengthSlice,
			mapEntry,
			scalarEntry,
			userDefinedScalarEntry,
		},
	}

	simpleTypeStringUnmarshaller = `package test

import (
	v5 "github.com/gofrs/uuid/v5"
	mint "github.com/vinyl-linux/mint"
	"io"
)

func (sf *SomeTestType) unmarshallSomeStringSlice(r io.Reader) (err error) {
	f := mint.NewSliceCollection(nil, false)
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	if f.Len() == 0 {
		sf.SomeStringSlice = nil
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
}
func (sf *SomeTestType) unmarshallSomeStringSlice(r io.Reader) (err error) {
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
}
func (sf *SomeTestType) unmarshallSomeStringSlice(r io.Reader) (err error) {
	f := mint.NewMapCollection(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer{})
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	if f.Len() == 0 {
		sf.SomeStringSlice = nil
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
}
func (sf *SomeTestType) unmarshallATypeOfSomeType(r io.Reader) (err error) {
	f := mint.NewUuidScalar(v5.UUID{})
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.ATypeOfSomeType = f.Value().(v5.UUID)
	return
}
func (sf *SomeTestType) unmarshallThingy(r io.Reader) (err error) {
	f := new(BlahType)
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Thingy = f.Value().(BlahType)
	return
}
func (sf *SomeTestType) Unmarshall(r io.Reader) (err error) {
	if err = sf.unmarshallSomeStringSlice(r); err != nil {
		return
	}
	if err = sf.unmarshallSomeStringSlice(r); err != nil {
		return
	}
	if err = sf.unmarshallSomeStringSlice(r); err != nil {
		return
	}
	if err = sf.unmarshallATypeOfSomeType(r); err != nil {
		return
	}
	if err = sf.unmarshallThingy(r); err != nil {
		return
	}
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	return
}
`

	simpleTypeStringMarshaller = `package test

import (
	mint "github.com/vinyl-linux/mint"
	"io"
)

func (sf SomeTestType) marshallSomeStringSlice(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.SomeStringSlice))
	for i := range f {
		f[i] = mint.NewStringScalar(sf.SomeStringSlice[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}
func (sf SomeTestType) marshallSomeStringSlice(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.SomeStringSlice))
	for i := range f {
		f[i] = mint.NewStringScalar(sf.SomeStringSlice[i])
	}
	return mint.NewSliceCollection(f, true).Marshall(w)
}
func (sf SomeTestType) marshallSomeStringSlice(w io.Writer) (err error) {
	f := make(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer)
	for k, v := range sf.SomeStringSlice {
		f[mint.NewStringScalar(k)] = mint.NewInt64Scalar(v)
	}
	return mint.NewMapCollection(f).Marshall(w)
}
func (sf SomeTestType) Marshall(w io.Writer) (err error) {
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	if err = sf.marshallSomeStringSlice(w); err != nil {
		return
	}
	if err = sf.marshallSomeStringSlice(w); err != nil {
		return
	}
	if err = sf.marshallSomeStringSlice(w); err != nil {
		return
	}
	if err = mint.NewUuidScalar(sf.ATypeOfSomeType).Marshall(w); err != nil {
		return
	}
	if err = sf.Thingy.Marshall(w); err != nil {
		return
	}
	return
}
`
)

func codeToString(c jen.Code) string {
	s := make(jen.Statement, 1)
	s[0] = c

	return s.GoString()
}

func codeSliceToFile(c []jen.Code) string {
	f := jen.NewFile("test")
	for _, stmnt := range c {
		f.Add(stmnt)
	}

	return f.GoString()
}
