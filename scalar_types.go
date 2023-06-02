package mint

import (
	"encoding/binary"
	"io"
	"time"
)

type StringScalar struct {
	v string
}

func NewStringScalar(s string) *StringScalar {
	return &StringScalar{
		v: s,
	}
}

func (s StringScalar) Marshall(w io.Writer) (err error) {
	err = binary.Write(w, binary.LittleEndian, int64(len(s.v)))
	if err != nil {
		return
	}

	_, err = w.Write([]byte(s.v))

	return
}

func (s *StringScalar) Unmarshall(r io.Reader) (err error) {
	var len int64
	err = binary.Read(r, binary.LittleEndian, &len)
	if err != nil {
		return
	}

	in := make([]byte, len)
	_, err = r.Read(in)
	if err != nil {
		return
	}

	s.v = string(in)

	return
}

func (s StringScalar) Value() any {
	return s.v
}

type DatetimeScalar struct {
	v time.Time
}

func NewDatetimeScalar(t time.Time) *DatetimeScalar {
	return &DatetimeScalar{
		v: t,
	}
}

func (s DatetimeScalar) Marshall(w io.Writer) error {
	intermediate := s.v.UnixNano()

	return binary.Write(w, binary.LittleEndian, intermediate)
}

func (s *DatetimeScalar) Unmarshall(r io.Reader) (err error) {
	var intermediate int64

	err = binary.Read(r, binary.LittleEndian, &intermediate)
	if err != nil {
		return
	}

	// This happens when an empty time.Time{} is serialised.
	//
	// In this situation, Unmarshall will create a time.Time with the date
	//   time.Date(1754, time.August, 30, 22, 42, 26, 128654848, time.Local)
	// Which represents the _actual_ earliest a time.Time can represent (as
	// opposed to 1st January, Year 0 which a nil time.Time seems to be)
	if intermediate == -6795364578871345152 {
		s.v = time.Time{}

		return
	}

	s.v = time.Unix(0, intermediate)

	return
}

func (s DatetimeScalar) Value() any {
	return s.v
}

type Int16Scalar struct {
	v int16
}

func NewInt16Scalar(i int16) *Int16Scalar {
	return &Int16Scalar{
		v: i,
	}
}

func (s *Int16Scalar) Marshall(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, s.v)
}

func (s *Int16Scalar) Unmarshall(r io.Reader) (err error) {
	return binary.Read(r, binary.LittleEndian, &s.v)
}

func (s Int16Scalar) Value() any {
	return s.v
}

type Int32Scalar struct {
	v int32
}

func NewInt32Scalar(i int32) *Int32Scalar {
	return &Int32Scalar{
		v: i,
	}
}

func (s *Int32Scalar) Marshall(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, s.v)
}

func (s *Int32Scalar) Unmarshall(r io.Reader) (err error) {
	return binary.Read(r, binary.LittleEndian, &s.v)
}

func (s Int32Scalar) Value() any {
	return s.v
}

type UInt32Scalar struct {
	v uint32
}

func NewUInt32Scalar(i uint32) *UInt32Scalar {
	return &UInt32Scalar{
		v: i,
	}
}

func (s *UInt32Scalar) Marshall(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, s.v)
}

func (s *UInt32Scalar) Unmarshall(r io.Reader) (err error) {
	return binary.Read(r, binary.LittleEndian, &s.v)
}

func (s UInt32Scalar) Value() any {
	return s.v
}

type Int64Scalar struct {
	v int64
}

func NewInt64Scalar(i int64) *Int64Scalar {
	return &Int64Scalar{
		v: i,
	}
}

func (s *Int64Scalar) Marshall(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, s.v)
}

func (s *Int64Scalar) Unmarshall(r io.Reader) (err error) {
	return binary.Read(r, binary.LittleEndian, &s.v)
}

func (s Int64Scalar) Value() any {
	return s.v
}

type Float32Scalar struct {
	v float32
}

func NewFloat32Scalar(i float32) *Float32Scalar {
	return &Float32Scalar{
		v: i,
	}
}

func (s *Float32Scalar) Marshall(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, s.v)
}

func (s *Float32Scalar) Unmarshall(r io.Reader) (err error) {
	return binary.Read(r, binary.LittleEndian, &s.v)
}

func (s Float32Scalar) Value() any {
	return s.v
}

type Float64Scalar struct {
	v float64
}

func NewFloat64Scalar(i float64) *Float64Scalar {
	return &Float64Scalar{
		v: i,
	}
}

func (s *Float64Scalar) Marshall(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, s.v)
}

func (s *Float64Scalar) Unmarshall(r io.Reader) (err error) {
	return binary.Read(r, binary.LittleEndian, &s.v)
}

func (s Float64Scalar) Value() any {
	return s.v
}

type ByteScalar struct {
	v byte
}

func NewByteScalar(i byte) *ByteScalar {
	return &ByteScalar{
		v: i,
	}
}

func (s *ByteScalar) Marshall(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, s.v)
}

func (s *ByteScalar) Unmarshall(r io.Reader) (err error) {
	return binary.Read(r, binary.LittleEndian, &s.v)
}

func (s ByteScalar) Value() any {
	return s.v
}

type BoolScalar struct {
	v bool
}

func NewBoolScalar(i bool) *BoolScalar {
	return &BoolScalar{
		v: i,
	}
}

func (s *BoolScalar) Marshall(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, s.v)
}

func (s *BoolScalar) Unmarshall(r io.Reader) (err error) {
	return binary.Read(r, binary.LittleEndian, &s.v)
}

func (s BoolScalar) Value() any {
	return s.v
}
