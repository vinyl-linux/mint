package parser

type scalar struct {
	size          int
	nulTerminated bool
}

var (
	scalars = map[string]scalar{
		"string":   scalar{0, true},
		"int":      scalar{32, false},
		"int32":    scalar{32, false},
		"int16":    scalar{16, false},
		"int64":    scalar{64, false},
		"float32":  scalar{32, false},
		"float64":  scalar{64, false},
		"byte":     scalar{8, false},
		"boolean":  scalar{8, false},
		"datetime": scalar{64, false},
	}
)
