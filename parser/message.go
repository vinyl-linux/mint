package parser

import (
	"fmt"
	"unicode"

	"github.com/alecthomas/participle/v2/lexer"
)

type AST struct {
	Types []AnnotatedType
	Enums []Enum
}

func toAST(d Document) (ad *AST, err error) {
	ad = new(AST)
	ad.Enums = make([]Enum, 0)
	ad.Types = make([]AnnotatedType, 0)

	for _, e := range d.Entries {
		if e.Enum != nil {
			ad.Enums = append(ad.Enums, *e.Enum)

			continue
		}

		at, err := toAnnotatedType(*e.Type)
		if err != nil {
			return nil, err
		}

		ad.Types = append(ad.Types, at)
	}

	return
}

type AnnotatedType struct {
	Pos     lexer.Position
	Name    string
	Entries []AnnotatedEntry
}

func (at AnnotatedType) name() string {
	return at.Name
}

func (at AnnotatedType) pos() lexer.Position {
	return at.Pos
}

type AnnotatedEntry struct {
	Field

	DocString       string
	Validations     []Validation
	Transformations []Transformation
}

func (ae *AnnotatedEntry) AppendDocString(s string) {
	if len(ae.DocString) == 0 {
		ae.DocString = s

		return
	}
	ae.DocString = ae.DocString + " " + s
}

func (ae *AnnotatedEntry) AppendValidation(v Validation) {
	if ae.Validations == nil {
		ae.Validations = make([]Validation, 0)
	}

	ae.Validations = append(ae.Validations, v)
}

func (ae *AnnotatedEntry) AppendTransformation(v Transformation) {
	if ae.Transformations == nil {
		ae.Transformations = make([]Transformation, 0)
	}

	ae.Transformations = append(ae.Transformations, v)
}

// IsValidType returns an error unless:
//
//  1. The specified type starts with a lower case and exists in our base scalars map; or
//  2. It starts with an upper case and exists as a Type or Enum in our AST
func (ae *AnnotatedEntry) IsValidType(names map[string][]lexer.Position) error {
	var (
		pos = ae.DataType.Pos

		t []string
	)

	switch {
	case ae.DataType.Scalar != nil:
		t = []string{ae.DataType.Scalar.Type}

	case ae.DataType.Slice != nil:
		t = []string{ae.DataType.Slice.Type}

	case ae.DataType.FixedSizeSlice != nil:
		t = []string{ae.DataType.FixedSizeSlice.Type}

	case ae.DataType.Map != nil:
		t = []string{ae.DataType.Map.Key, ae.DataType.Map.Value}

	default:
		return fmt.Errorf("Unable to determine type at %s", pos.String())
	}

	for _, s := range t {
		if !scalarOrNames(s, names) {
			return incorrectTypeErr{
				t:   s,
				pos: ae.DataType.Pos,
			}
		}

	}

	return nil
}

type Validation struct {
	IsCustom bool
	Function string
}

type Transformation struct {
	IsCustom bool
	Function string
}

type incorrectTypeErr struct {
	t   string
	pos lexer.Position
}

func (e incorrectTypeErr) Error() string {
	return fmt.Sprintf("Unrecognised data type %s at %s", e.t, e.pos.String())
}

func scalarOrNames(s string, names map[string][]lexer.Position) bool {
	if unicode.IsLower(rune(s[0])) {
		_, ok := Scalars[s]
		return ok
	}

	_, ok := names[s]
	return ok
}
