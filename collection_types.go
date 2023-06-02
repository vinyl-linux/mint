package mint

import (
	"encoding/binary"
	"io"
)

type SliceCollection struct {
	fixedLength bool
	len         uint32
	V           []MarshallerUnmarshallerValuer
}

func NewSliceCollection(v []MarshallerUnmarshallerValuer, isFixedLength bool) *SliceCollection {
	return &SliceCollection{
		fixedLength: isFixedLength,
		len:         uint32(len(v)),
		V:           v,
	}
}

func (s SliceCollection) Len() int {
	return int(s.len)
}

func (s *SliceCollection) ReadSize(r io.Reader) (err error) {
	if s.fixedLength {
		return nil
	}

	return binary.Read(r, binary.LittleEndian, &s.len)
}

func (s SliceCollection) Marshall(w io.Writer) (err error) {
	if !s.fixedLength {
		err = binary.Write(w, binary.LittleEndian, s.len)
		if err != nil {
			return
		}
	}

	for _, i := range s.V {
		err = i.Marshall(w)
		if err != nil {
			return
		}
	}

	return nil
}

func (s *SliceCollection) Unmarshall(r io.Reader) (err error) {
	for i := uint32(0); i < s.len; i++ {
		err = s.V[i].Unmarshall(r)
		if err != nil {
			return
		}
	}

	return
}

func (s SliceCollection) Value() any {
	return s.V
}

type MapCollection struct {
	V   map[MarshallerUnmarshallerValuer]MarshallerUnmarshallerValuer
	len uint32
}

func NewMapCollection(v map[MarshallerUnmarshallerValuer]MarshallerUnmarshallerValuer) *MapCollection {
	return &MapCollection{
		V: v,
	}
}

func (s MapCollection) Len() int {
	return int(s.len)
}

func (s *MapCollection) ReadSize(r io.Reader) (err error) {
	var intermediate uint32 = 0

	err = binary.Read(r, binary.LittleEndian, &intermediate)
	if err != nil {
		return
	}

	s.len = intermediate / 2

	return
}

func (s MapCollection) Marshall(w io.Writer) (err error) {
	return NewSliceCollection(s.slice(), false).Marshall(w)
}

func (s *MapCollection) Unmarshall(r io.Reader) (err error) {
	sl := s.slice()

	err = NewSliceCollection(sl, false).Unmarshall(r)
	if err != nil {
		return
	}

	s.V = make(map[MarshallerUnmarshallerValuer]MarshallerUnmarshallerValuer)

	for i := 0; i < len(sl); i += 2 {
		s.V[sl[i]] = sl[i+1]

	}

	return
}

func (s MapCollection) Value() any {
	return s.V
}

func (s MapCollection) slice() []MarshallerUnmarshallerValuer {
	slice := make([]MarshallerUnmarshallerValuer, len(s.V)*2)

	idx := 0
	for k, v := range s.V {
		slice[idx] = k
		slice[idx+1] = v

		idx += 2
	}

	return slice
}
