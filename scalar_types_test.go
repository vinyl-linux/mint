package mint

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestStringScalar(t *testing.T) {
	for _, test := range []struct {
		name           string
		s              string
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Empty string should return empty string", "", false, false},
		{"Arbitrary string should return same string", "Hello, World!", false, false},
		{"Massive string should return exactly", makeLongString(), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewStringScalar(test.s)
			if v == nil {
				t.Fatalf("expected instance of StringScala, received null")
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
				v = NewStringScalar("")

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
				received := v.Value().(string)
				if test.s != received {
					t.Errorf("expected %#v, received %#v", test.s, received)
				}
			})
		})
	}
}

func TestDatetimeScalar(t *testing.T) {
	for _, test := range []struct {
		name           string
		d              time.Time
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Empty time.Time should return empty time.Time", time.Time{}, false, false},
		{"Specific time.Time should return same", time.Now(), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewDatetimeScalar(test.d)
			if v == nil {
				t.Fatalf("expected instance of TimeDateScalar, received null")
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
				v = new(DatetimeScalar)

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
				received := v.Value().(time.Time)
				if !test.d.Equal(received) {
					t.Errorf("expected %#v, received %#v", test.d, received)
				}
			})
		})
	}
}

func TestInt16Scalar(t *testing.T) {
	for _, test := range []struct {
		name           string
		i              int16
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Zero value should remain as such", int16(0), false, false},
		{"Malformed zero value should remain as such", int16(0000000), false, false},
		{"Large number with underscore delimeter should remain as such", int16(10_000), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewInt16Scalar(test.i)
			if v == nil {
				t.Fatalf("expected instance of Int16Scalar, received null")
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
				v = new(Int16Scalar)

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
				received := v.Value().(int16)
				if test.i != received {
					t.Errorf("expected %v, received %v", test.i, received)
				}
			})
		})
	}
}

func TestInt32Scalar(t *testing.T) {
	for _, test := range []struct {
		name           string
		i              int32
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Zero value should remain as such", int32(0), false, false},
		{"Malformed zero value should remain as such", int32(0000000), false, false},
		{"Large number with underscore delimeter should remain as such", int32(10_000), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewInt32Scalar(test.i)
			if v == nil {
				t.Fatalf("expected instance of Int32Scalar, received null")
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
				v = new(Int32Scalar)

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
				received := v.Value().(int32)
				if test.i != received {
					t.Errorf("expected %v, received %v", test.i, received)
				}
			})
		})
	}
}

func TestUInt32Scalar(t *testing.T) {
	for _, test := range []struct {
		name           string
		i              uint32
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Zero value should remain as such", uint32(0), false, false},
		{"Malformed zero value should remain as such", uint32(0000000), false, false},
		{"Large number with underscore delimeter should remain as such", uint32(10_000), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewUInt32Scalar(test.i)
			if v == nil {
				t.Fatalf("expected instance of Unint32Scalar, received null")
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
				v = new(UInt32Scalar)

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
				received := v.Value().(uint32)
				if test.i != received {
					t.Errorf("expected %v, received %v", test.i, received)
				}
			})
		})
	}
}

func TestInt64Scalar(t *testing.T) {
	for _, test := range []struct {
		name           string
		i              int64
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Zero value should remain as such", int64(0), false, false},
		{"Malformed zero value should remain as such", int64(0000000), false, false},
		{"Large number with underscore delimeter should remain as such", int64(10_000), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewInt64Scalar(test.i)
			if v == nil {
				t.Fatalf("expected instance of Int64Scalar, received null")
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
				v = new(Int64Scalar)

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
				received := v.Value().(int64)
				if test.i != received {
					t.Errorf("expected %v, received %v", test.i, received)
				}
			})
		})
	}
}

func TestFloat32Scalar(t *testing.T) {
	for _, test := range []struct {
		name           string
		i              float32
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Zero value should remain as such", float32(0), false, false},
		{"Malformed zero value should remain as such", float32(0000000), false, false},
		{"Large number with underscore delimeter should remain as such", float32(10_000), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewFloat32Scalar(test.i)
			if v == nil {
				t.Fatalf("expected instance of Float32Scalar, received null")
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
				v = new(Float32Scalar)

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
				received := v.Value().(float32)
				if test.i != received {
					t.Errorf("expected %v, received %v", test.i, received)
				}
			})
		})
	}
}

func TestFloat64Scalar(t *testing.T) {
	for _, test := range []struct {
		name           string
		i              float64
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Zero value should remain as such", float64(0), false, false},
		{"Malformed zero value should remain as such", float64(0000000), false, false},
		{"Large number with underscore delimeter should remain as such", float64(10_000), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewFloat64Scalar(test.i)
			if v == nil {
				t.Fatalf("expected instance of Float64Scalar, received null")
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
				v = new(Float64Scalar)

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
				received := v.Value().(float64)
				if test.i != received {
					t.Errorf("expected %v, received %v", test.i, received)
				}
			})
		})
	}
}

func TestByteScalar(t *testing.T) {
	for _, test := range []struct {
		name           string
		i              byte
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Zero value should remain as such", byte(0), false, false},
		{"Malformed zero value should remain as such", byte(0000000), false, false},
		{"Character byte should remain as such", byte('a'), false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewByteScalar(test.i)
			if v == nil {
				t.Fatalf("expected instance of ByteScalar, received null")
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
				v = new(ByteScalar)

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
				received := v.Value().(byte)
				if test.i != received {
					t.Errorf("expected %v, received %v", test.i, received)
				}
			})
		})
	}
}

func TestBoolScalar(t *testing.T) {
	var zeroedBool bool

	for _, test := range []struct {
		name           string
		i              bool
		expectWriteErr bool
		expectReadErr  bool
	}{
		{"Zero value should remain as such", zeroedBool, false, false},
		{"False value should remain as such", false, false, false},
		{"True value should remain as such", true, false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			v := NewBoolScalar(test.i)
			if v == nil {
				t.Fatalf("expected instance of BoolScalar, received null")
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
				v = new(BoolScalar)

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
				received := v.Value().(bool)
				if test.i != received {
					t.Errorf("expected %v, received %v", test.i, received)
				}
			})
		})
	}
}

func makeLongString() string {
	sb := strings.Builder{}
	for i := 0; i < 100_000; i++ {
		sb.WriteString("hello")
	}

	return sb.String()
}
