package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestParse(t *testing.T) {
	for _, test := range []struct {
		name      string
		body      string
		expect    *AST
		expectErr bool
	}{
		{"empty body returns empty ast", "", emptyAST, false},
		{"colliding type names error", collidingNames, nil, true},
		{"colliding field names error", collidingFieldNames, nil, true},
		{"colliding field tags error", collidingTags, nil, true},
		{"missing field tags error", missingFields, nil, true},
		{"valid input works", validInput, fullAST, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			received, err := Parse(test.name, strings.NewReader(test.body))
			if err != nil {
				t.Log(err)
			}

			if err == nil && test.expectErr {
				t.Errorf("expected error, received none")
			} else if err != nil && !test.expectErr {
				t.Errorf("unexpected error %#v", err)
			}

			if !reflect.DeepEqual(test.expect, received) {
				receivedStr := repr.String(received)
				expectStr := repr.String(test.expect)

				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(expectStr, receivedStr, false)

				t.Error(dmp.DiffPrettyText(diffs))
			}
		})
	}
}

var (
	emptyAST = new(AST)
	fullAST  = &AST{
		Types: []annotatedType{
			annotatedType{
				Pos: lexer.Position{
					Filename: "valid input works",
					Offset:   1,
					Line:     2,
					Column:   1,
				},
				Name: "Foo",
				Entries: []annotatedEntry{
					annotatedEntry{
						Field: Field{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   178,
								Line:     8,
								Column:   3,
							},
							DataType: &DataType{
								Pos: lexer.Position{
									Filename: "valid input works",
									Offset:   178,
									Line:     8,
									Column:   3,
								},
								Scalar: String,
							},
							Name: "Hello",
							Tag:  0,
						},
						DocString: "Hello, world! This is a multiline doc string :)",
						Validations: []validation{
							{
								Function: "not_empty",
							},
						},
						Transformations: []transformation{
							{
								Function: "to_lowercase",
							},
							{
								IsCustom: true,
								Function: "to_korean",
							},
						},
					},
				},
			},
		},
		Enums: []Enum{
			{
				Pos: lexer.Position{
					Filename: "valid input works",
					Offset:   199,
					Line:     11,
					Column:   1,
				},
				Name: "Bar",
				Values: []*EnumEntry{
					{
						Pos: lexer.Position{
							Filename: "valid input works",
							Offset:   212,
							Line:     12,
							Column:   3,
						},
						Value: &EnumValue{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   212,
								Line:     12,
								Column:   3,
							},
							Key:   "A",
							Value: 1,
						},
					},
					{
						Pos: lexer.Position{
							Filename: "valid input works",
							Offset:   221,
							Line:     13,
							Column:   3,
						},
						Value: &EnumValue{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   221,
								Line:     13,
								Column:   3,
							},
							Key:   "B",
							Value: 2},
					},
					{
						Pos: lexer.Position{
							Filename: "valid input works",
							Offset:   230,
							Line:     14,
							Column:   3,
						},
						Value: &EnumValue{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   230,
								Line:     14,
								Column:   3,
							},
							Key:   "C",
							Value: 3,
						},
					},
				},
			},
		},
	}

	collidingNames = `
type Foo {}
type Foo {}
type Bar {}
`
	collidingFieldNames = `
type Foo {
  string Bar = 0;
  string Bar = 1;
}
`
	collidingTags = `
type Foo {
  string Bar = 0;
  string Baz = 0;
}
`
	missingFields = `
type Foo {
  string Bar = 0;
  string Baz = 2;
}
`
	validInput = `
type Foo {
  +mint:doc:"Hello, world!"
  +mint:doc:"This is a multiline doc string :)"
  +mint:validate:not_empty
  +mint:transform:to_lowercase
  +custom:transform:to_korean
  string Hello = 0;
}

enum Bar {
  A = 1;
  B = 2;
  C = 3;
}
`
)
