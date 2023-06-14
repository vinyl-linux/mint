package benchmarks

import (
	v5 "github.com/gofrs/uuid/v5"
	mint "github.com/vinyl-linux/mint"
	"io"
)

type Benchmarker struct {
	ID               v5.UUID
	ShortString      string
	LongString       string
	ManyShortStrings []string
	ManyLongStrings  []string
	SomeNumber       int64
}

func (sf Benchmarker) Validate() error {
	errors := make([]error, 0)
	for _, err := range []error{} {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return mint.ValidationErrors("Benchmarker", errors)
}
func (sf *Benchmarker) Transform() (err error) {
	return
}
func (sf Benchmarker) Value() any {
	return sf
}
func (sf *Benchmarker) unmarshallID(r io.Reader) (err error) {
	f := mint.NewUuidScalar(v5.UUID{})
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.ID = f.Value().(v5.UUID)
	return
}
func (sf *Benchmarker) unmarshallShortString(r io.Reader) (err error) {
	f := mint.NewStringScalar("")
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.ShortString = f.Value().(string)
	return
}
func (sf *Benchmarker) unmarshallLongString(r io.Reader) (err error) {
	f := mint.NewStringScalar("")
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.LongString = f.Value().(string)
	return
}
func (sf *Benchmarker) unmarshallManyShortStrings(r io.Reader) (err error) {
	f := mint.NewSliceCollection(nil, false)
	f.ReadSize(r)
	f.V = make([]mint.MarshallerUnmarshallerValuer, f.Len())
	for i := range f.V {
		f.V[i] = mint.NewStringScalar("")
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.ManyShortStrings = make([]string, f.Len())
	for i, v := range f.Value().([]mint.MarshallerUnmarshallerValuer) {
		sf.ManyShortStrings[i] = v.Value().(string)
	}
	return
}
func (sf *Benchmarker) unmarshallManyLongStrings(r io.Reader) (err error) {
	f := mint.NewSliceCollection(nil, false)
	f.ReadSize(r)
	f.V = make([]mint.MarshallerUnmarshallerValuer, f.Len())
	for i := range f.V {
		f.V[i] = mint.NewStringScalar("")
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.ManyLongStrings = make([]string, f.Len())
	for i, v := range f.Value().([]mint.MarshallerUnmarshallerValuer) {
		sf.ManyLongStrings[i] = v.Value().(string)
	}
	return
}
func (sf *Benchmarker) unmarshallSomeNumber(r io.Reader) (err error) {
	f := mint.NewInt64Scalar(int64(0))
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.SomeNumber = f.Value().(int64)
	return
}
func (sf *Benchmarker) Unmarshall(r io.Reader) (err error) {
	if err = sf.unmarshallID(r); err != nil {
		return
	}
	if err = sf.unmarshallShortString(r); err != nil {
		return
	}
	if err = sf.unmarshallLongString(r); err != nil {
		return
	}
	if err = sf.unmarshallManyShortStrings(r); err != nil {
		return
	}
	if err = sf.unmarshallManyLongStrings(r); err != nil {
		return
	}
	if err = sf.unmarshallSomeNumber(r); err != nil {
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
func (sf Benchmarker) marshallManyShortStrings(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.ManyShortStrings))
	for i := range f {
		f[i] = mint.NewStringScalar(sf.ManyShortStrings[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}
func (sf Benchmarker) marshallManyLongStrings(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.ManyLongStrings))
	for i := range f {
		f[i] = mint.NewStringScalar(sf.ManyLongStrings[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}
func (sf Benchmarker) Marshall(w io.Writer) (err error) {
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	if err = mint.NewUuidScalar(sf.ID).Marshall(w); err != nil {
		return
	}
	if err = mint.NewStringScalar(sf.ShortString).Marshall(w); err != nil {
		return
	}
	if err = mint.NewStringScalar(sf.LongString).Marshall(w); err != nil {
		return
	}
	if err = sf.marshallManyShortStrings(w); err != nil {
		return
	}
	if err = sf.marshallManyLongStrings(w); err != nil {
		return
	}
	if err = mint.NewInt64Scalar(sf.SomeNumber).Marshall(w); err != nil {
		return
	}
	return
}
