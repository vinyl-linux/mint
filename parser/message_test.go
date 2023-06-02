package parser

import (
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
)

func TestAnnotatedEntry_IsValidType(t *testing.T) {
	ae := &annotatedEntry{
		Field: Field{
			DataType: &DataType{},
		},
	}

	err := ae.IsValidType(map[string][]lexer.Position{})
	if err == nil {
		t.Errorf("expected error, got none")
	}
}
