package parser

var (
	scalars = map[string]bool{
		"string":   false,
		"datetime": false,

		"int16":   false,
		"int32":   false,
		"int64":   false,
		"float32": false,
		"float64": false,
		"byte":    false,
		"bool":    false,
	}
)
