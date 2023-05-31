package mint

import (
	"io"
)

type Marshaller interface {
	Marshall(io.Writer) error
}

type Unmarshaller interface {
	Unmarshall(io.Reader) error
}

type Valuer interface {
	Value() any
}

type MarshallerUnmarshallerValuer interface {
	Marshaller
	Unmarshaller
	Valuer
}
