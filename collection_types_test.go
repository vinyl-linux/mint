package mint

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
)

type customMUV struct {
	foo string
	bar int64
	baz bool
}

func (c customMUV) Marshall(w io.Writer) (err error) {
	_ = NewStringScalar(c.foo).Marshall(w)
	_ = NewInt64Scalar(c.bar).Marshall(w)
	_ = NewBoolScalar(c.baz).Marshall(w)

	return
}

func (c *customMUV) Unmarshall(r io.Reader) (err error) {
	cs := NewStringScalar("")
	cs.Unmarshall(r)
	c.foo = cs.Value().(string)

	ci := NewInt64Scalar(0)
	ci.Unmarshall(r)
	c.bar = ci.Value().(int64)

	cb := NewBoolScalar(false)
	cb.Unmarshall(r)
	c.baz = cb.Value().(bool)

	return
}

func (c customMUV) Value() any {
	return c
}

func TestSliceCollection_Scalars(t *testing.T) {
	for _, test := range []struct {
		name           string
		s              []string
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Empty string slice should return empty string", []string{""}, false, false},
		{"Single length slice of arbitrary string data should return exactly", []string{"Hello, World!"}, false, false},
		{"Single length slice, massive string should return exactly", []string{makeLongString()}, false, false},
		{"Massive slice of strings should return exactly", makeLongStringSlice(), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			ss := make([]MarshallerUnmarshallerValuer, len(test.s))
			for idx, elem := range test.s {
				ss[idx] = NewStringScalar(elem)
			}

			v := NewSliceCollection(ss, false)
			if v == nil {
				t.Fatalf("expected instance of SliceCollection, received null")
			}

			b := new(bytes.Buffer)
			t.Run("Marshall", func(t *testing.T) {
				err := v.Marshall(b)
				if err == nil && test.expectWriteErr {
					t.Errorf("expected error, received none")
				} else if err != nil && !test.expectWriteErr {
					t.Errorf("unexpected error %#v", err)
				}

				if b.Len() == 0 {
					t.Errorf("no bytes were written")
				}
			})

			t.Run("Unmarshall", func(t *testing.T) {
				v = NewSliceCollection(nil, false)
				v.ReadSize(b)

				v.V = make([]MarshallerUnmarshallerValuer, v.Len())
				for idx := range v.V {
					v.V[idx] = NewStringScalar("")
				}

				err := v.Unmarshall(b)
				if err == nil && test.expectReadErr {
					t.Errorf("expected error, received none")
				} else if err != nil && !test.expectReadErr {
					t.Errorf("unexpected error %#v", err)
				}

				if b.Len() > 0 {
					t.Errorf("bytes left over: %#v", b.Bytes())
				}
			})

			t.Run("Value", func(t *testing.T) {
				received := make([]string, v.Len())
				for idx, elem := range v.Value().([]MarshallerUnmarshallerValuer) {
					str, _ := elem.Value().(string)
					received[idx] = str
				}

				if !reflect.DeepEqual(test.s, received) {
					t.Errorf("expected %#v, received %#v", test.s, received)
				}
			})
		})
	}
}

func TestSliceCollection_FixedSize(t *testing.T) {
	expect := make([]bool, 5)
	for idx := range expect {
		expect[idx] = true
	}

	ss := make([]MarshallerUnmarshallerValuer, 5)
	for idx := range ss {
		ss[idx] = NewBoolScalar(expect[idx])
	}

	in := NewSliceCollection(ss, true)
	b := new(bytes.Buffer)

	err := in.Marshall(b)
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	ov := make([]MarshallerUnmarshallerValuer, 5)
	for idx := range ov {
		ov[idx] = NewBoolScalar(false)
	}
	out := NewSliceCollection(ov, true)
	out.ReadSize(b)

	err = out.Unmarshall(b)
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	received := make([]bool, out.Len())
	for idx, elem := range out.Value().([]MarshallerUnmarshallerValuer) {
		received[idx] = elem.Value().(bool)
	}

	if !reflect.DeepEqual(expect, received) {
		t.Errorf("expected\n\t%#v\nreceived\n\t%#v", expect, received)
	}
}

func TestSliceCollection_Complex(t *testing.T) {
	expect := make([]*customMUV, 5)
	for idx := range expect {
		expect[idx] = &customMUV{
			foo: "hello, world!",
			bar: 12345,
			baz: true,
		}
	}

	ss := make([]MarshallerUnmarshallerValuer, 5)
	for idx := range ss {
		ss[idx] = expect[idx]
	}

	in := NewSliceCollection(ss, true)
	b := new(bytes.Buffer)

	err := in.Marshall(b)
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	ov := make([]MarshallerUnmarshallerValuer, 5)
	for idx := range ov {
		ov[idx] = new(customMUV)
	}
	out := NewSliceCollection(ov, true)

	err = out.Unmarshall(b)
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	received := make([]*customMUV, out.Len())
	for idx, elem := range out.Value().([]MarshallerUnmarshallerValuer) {
		received[idx] = elem.(*customMUV)
	}

	if !reflect.DeepEqual(expect, received) {
		t.Errorf("expected\n\t%#v\nreceived\n\t%#v", expect, received)
	}
}

func TestMapCollection(t *testing.T) {
	for _, test := range []struct {
		name           string
		m              map[string]bool
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Empty map", map[string]bool{}, false, false},
		{"Single entry map", map[string]bool{"Hello": true}, false, false},
		{"Large map", makeLongStringBoolMap(), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			ss := make(map[MarshallerUnmarshallerValuer]MarshallerUnmarshallerValuer)
			for k, v := range test.m {
				ss[NewStringScalar(k)] = NewBoolScalar(v)
			}

			v := NewMapCollection(ss)
			if v == nil {
				t.Fatalf("expected instance of SliceCollection, received null")
			}

			b := new(bytes.Buffer)
			t.Run("Marshall", func(t *testing.T) {
				err := v.Marshall(b)
				if err == nil && test.expectWriteErr {
					t.Errorf("expected error, received none")
				} else if err != nil && !test.expectWriteErr {
					t.Errorf("unexpected error %#v", err)
				}

				if b.Len() == 0 {
					t.Errorf("no bytes were written")
				}
			})

			t.Run("Unmarshall", func(t *testing.T) {
				v = NewMapCollection(nil)
				v.ReadSize(b)

				v.V = make(map[MarshallerUnmarshallerValuer]MarshallerUnmarshallerValuer)
				for i := 0; i < v.Len(); i++ {
					v.V[NewStringScalar("")] = NewBoolScalar(false)
				}

				err := v.Unmarshall(b)
				if err == nil && test.expectReadErr {
					t.Errorf("expected error, received none")
				} else if err != nil && !test.expectReadErr {
					t.Errorf("unexpected error %#v", err)
				}

				if b.Len() > 0 {
					t.Errorf("bytes left over: %#v", b.Bytes())
				}
			})

			t.Run("Value", func(t *testing.T) {
				received := make(map[string]bool, v.Len())
				for k, v := range v.Value().(map[MarshallerUnmarshallerValuer]MarshallerUnmarshallerValuer) {
					kS, _ := k.Value().(string)
					vS, _ := v.Value().(bool)

					received[kS] = vS
				}

				if !reflect.DeepEqual(test.m, received) {
					t.Errorf("expected %#v, received %#v", test.m, received)
				}
			})
		})
	}
}

func makeLongStringSlice() (out []string) {
	out = make([]string, 100_000)
	for i := range out {
		out[i] = "hello"
	}

	return
}

func makeLongStringBoolMap() (out map[string]bool) {
	out = make(map[string]bool)

	for i := 0; i < 1000; i++ {
		out[fmt.Sprintf("hello %d", i)] = true
	}

	return
}
