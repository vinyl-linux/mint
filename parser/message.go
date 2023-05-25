package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type AST struct {
	Types []annotatedType
	Enums []Enum
}

func toAST(d Document) (ad *AST, err error) {
	ad = new(AST)
	ad.Enums = make([]Enum, 0)
	ad.Types = make([]annotatedType, 0)

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

type annotatedType struct {
	Pos     lexer.Position
	Name    string
	Entries []annotatedEntry
}

func (at annotatedType) name() string {
	return at.Name
}

func (at annotatedType) pos() lexer.Position {
	return at.Pos
}

type annotatedEntry struct {
	Field

	DocString       string
	Validations     []validation
	Transformations []transformation
}

func (ae annotatedEntry) name() string {
	return ae.Name
}

func (ae annotatedEntry) pos() lexer.Position {
	return ae.Pos
}

func (ae *annotatedEntry) AppendDocString(s string) {
	if len(ae.DocString) == 0 {
		ae.DocString = s

		return
	}
	ae.DocString = ae.DocString + " " + s
}

func (ae *annotatedEntry) AppendValidation(v validation) {
	if ae.Validations == nil {
		ae.Validations = make([]validation, 0)
	}

	ae.Validations = append(ae.Validations, v)
}

func (ae *annotatedEntry) AppendTransformation(v transformation) {
	if ae.Transformations == nil {
		ae.Transformations = make([]transformation, 0)
	}

	ae.Transformations = append(ae.Transformations, v)
}

type validation struct {
	IsCustom bool
	Function string
}

type transformation struct {
	IsCustom bool
	Function string
}
