package parser

var (
	scalars = map[string]bool{
		"string":   true,
		"datetime": true,
		"uuid":     true,

		"int16":   true,
		"int32":   true,
		"int64":   true,
		"float32": true,
		"float64": true,
		"byte":    true,
		"bool":    true,
	}
)
