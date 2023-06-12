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
		{"invalid input errors", "hello, world!", nil, true},
		{"empty body returns empty ast", "", emptyAST, false},
		{"colliding type names error", collidingNames, nil, true},
		{"colliding field names error", collidingFieldNames, nil, true},
		{"colliding field tags error", collidingTags, nil, true},
		{"missing field tags error", missingFields, nil, true},
		{"valid input works", validInput, fullAST, false},
		{"invalid scalar type errors", invalidScalar, nil, true},
		{"invalid map key type errors", invalidMapKey, nil, true},
		{"invalid map value type errors", invalidMapValue, nil, true},
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

func TestParseFile(t *testing.T) {
	for _, test := range []struct {
		fn        string
		expectErr bool
	}{
		{"testdata/nonsuch.mint", true},
		{"testdata/valid-documents/location.mint", false},
	} {
		t.Run(test.fn, func(t *testing.T) {
			_, err := ParseFile(test.fn)
			if err == nil && test.expectErr {
				t.Errorf("expected error, received none")
			} else if err != nil && !test.expectErr {
				t.Errorf("unexpected error %s", err)
			}
		})
	}
}

func TestParseDir(t *testing.T) {
	for _, test := range []struct {
		dir       string
		expectErr bool
	}{
		{"testdata/valid-documents", false},
		{"testdata/nonsuch", true},
		{"testdata/invalid-document", true},
	} {
		t.Run(test.dir, func(t *testing.T) {
			_, err := ParseDir(test.dir)
			if err == nil && test.expectErr {
				t.Errorf("expected error, received none")
			} else if err != nil && !test.expectErr {
				t.Errorf("unexpected error %s", err)
			}
		})
	}
}

var (
	emptyAST = new(AST)
	fullAST  = &AST{
		Types: []AnnotatedType{
			AnnotatedType{
				Pos: lexer.Position{
					Filename: "valid input works",
					Offset:   1,
					Line:     2,
					Column:   1,
				},
				Name: "Foo",
				Entries: []AnnotatedEntry{
					AnnotatedEntry{
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
								Scalar: &Scalar{
									Pos: lexer.Position{
										Filename: "valid input works",
										Offset:   178,
										Line:     8,
										Column:   3,
									},
									Type: "string",
								},
							},
							Name: "Hello",
							Tag:  0,
						},
						DocString: "Hello, world! This is a multiline doc string :)",
						Validations: []Validation{
							{
								Function: "not_empty",
							},
						},
						Transformations: []Transformation{
							{
								Function: "to_lowercase",
							},
							{
								IsCustom: true,
								Function: "to_korean",
							},
						},
					},
					{
						Field: Field{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   228,
								Line:     11,
								Column:   3,
							},
							DataType: &DataType{
								Pos: lexer.Position{
									Filename: "valid input works",
									Offset:   228,
									Line:     11,
									Column:   3,
								},
								Scalar: &Scalar{
									Pos: lexer.Position{
										Filename: "valid input works",
										Offset:   228,
										Line:     11,
										Column:   3,
									},
									Type: "Bar",
								},
							},
							Name: "Bar",
							Tag:  1,
						},
						DocString: "Some bar value",
					},
					{
						Field: Field{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   244,
								Line:     13,
								Column:   3,
							},
							DataType: &DataType{
								Pos: lexer.Position{
									Filename: "valid input works",
									Offset:   244,
									Line:     13,
									Column:   3,
								},
								Map: &MapType{
									Pos: lexer.Position{
										Filename: "valid input works",
										Offset:   244,
										Line:     13,
										Column:   3,
									},
									Key:   "string",
									Value: "int32",
								},
							},
							Name: "MappyMap",
							Tag:  2,
						},
					},
					{
						Field: Field{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   279,
								Line:     14,
								Column:   3,
							},
							DataType: &DataType{
								Pos: lexer.Position{
									Filename: "valid input works",
									Offset:   279,
									Line:     14,
									Column:   3,
								},
								Slice: &SliceType{
									Pos: lexer.Position{
										Filename: "valid input works",
										Offset:   279,
										Line:     14,
										Column:   3,
									},
									Type: "string",
								},
							},
							Name: "UnboundString",
							Tag:  3,
						},
					},
					{
						Field: Field{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   309,
								Line:     15,
								Column:   3,
							},
							DataType: &DataType{
								Pos: lexer.Position{
									Filename: "valid input works",
									Offset:   309,
									Line:     15,
									Column:   3,
								},
								FixedSizeSlice: &FixedSizedSliceType{
									Pos: lexer.Position{
										Filename: "valid input works",
										Offset:   309,
										Line:     15,
										Column:   3,
									},
									Size: 5,
									Type: "byte",
								},
							},
							Name: "FixedSizeByteSlice",
							Tag:  4,
						},
					},
				},
			},
		},
		Enums: []Enum{
			{
				Pos: lexer.Position{
					Filename: "valid input works",
					Offset:   344,
					Line:     18,
					Column:   1,
				},
				Name: "Bar",
				Values: []*EnumEntry{
					{
						Pos: lexer.Position{
							Filename: "valid input works",
							Offset:   357,
							Line:     19,
							Column:   3,
						},
						Value: &EnumValue{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   357,
								Line:     19,
								Column:   3,
							},
							Key:   "A",
							Value: 1,
						},
					},
					{
						Pos: lexer.Position{
							Filename: "valid input works",
							Offset:   366,
							Line:     20,
							Column:   3,
						},
						Value: &EnumValue{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   366,
								Line:     20,
								Column:   3,
							},
							Key:   "B",
							Value: 2,
						},
					},
					{
						Pos: lexer.Position{
							Filename: "valid input works",
							Offset:   375,
							Line:     21,
							Column:   3,
						},
						Value: &EnumValue{
							Pos: lexer.Position{
								Filename: "valid input works",
								Offset:   375,
								Line:     21,
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

  +mint:doc:"Some bar value"
  Bar Bar = 1;

  map<string, int32> MappyMap = 2;
  []string UnboundString = 3;
  [5]byte FixedSizeByteSlice = 4;
}

enum Bar {
  A = 1;
  B = 2;
  C = 3;
}
`

	invalidScalar = `
type Foo {
  str Bar = 0;
}
`

	invalidMapKey = `
type Foo {
  map<Bar, string> BarMap = 0;
}
`
	invalidMapValue = `
type Foo {
  map<string, Bar> BarMap = 0;
}
`
)
