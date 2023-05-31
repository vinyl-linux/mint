package mint

import (
	"encoding/binary"
	"io"
)

type SliceCollection struct {
	fixedLength bool
	len         uint32
	v           []MarshallerUnmarshallerValuer
}

func NewSliceCollection(v []MarshallerUnmarshallerValuer, isFixedLength bool) *SliceCollection {
	return &SliceCollection{
		fixedLength: isFixedLength,
		len:         uint32(len(v)),
		v:           v,
	}
}

func (s SliceCollection) Marshall(w io.Writer) (err error) {
	if !s.fixedLength {
		binary.Write(w, binary.LittleEndian, s.len)
	}

	for _, i := range s.v {
		err = i.Marshall(w)
		if err != nil {
			return
		}
	}

	return nil
}

func (s *SliceCollection) Unmarshall(r io.Reader) (err error) {
	if s.fixedLength {
		binary.Read(r, binary.LittleEndian, &s.len)
	}

	for i := uint32(0); i < s.len; i++ {
		err = s.v[i].Unmarshall(r)
		if err != nil {
			return
		}
	}

	return
}

func (s SliceCollection) Value() any {
	return s.v
}
