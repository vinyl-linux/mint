package parser

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Document struct {
	Pos lexer.Position

	Entries []*Entry `( @@ ";"* )*`
}

type Entry struct {
	Pos lexer.Position

	Message *Message ` @@`
	Enum    *Enum    `| @@`
}

type Value struct {
	Pos lexer.Position

	String    *string  `  @String`
	Number    *float64 `| @Float`
	Int       *int64   `| @Int`
	Bool      *bool    `| (@"true" | "false")`
	Reference *string  `| @Ident @( "." Ident )*`
	Map       *Map     `| @@`
	Array     *Array   `| @@`
}

type Array struct {
	Pos lexer.Position

	Elements []*Value `"[" ( @@ ( ","? @@ )* )? "]"`
}

type Map struct {
	Pos lexer.Position

	Entries []*MapEntry `"{" ( @@ ( ( "," )? @@ )* )? "}"`
}

type MapEntry struct {
	Pos lexer.Position

	Key   *Value `@@`
	Value *Value `":"? @@`
}

type Enum struct {
	Pos lexer.Position

	Name   string       `"enum" @Ident`
	Values []*EnumEntry `"{" ( @@ ( ";" )* )* "}"`
}

type EnumEntry struct {
	Pos lexer.Position

	Value *EnumValue `  @@`
}

type EnumValue struct {
	Pos lexer.Position

	Key   string `@Ident`
	Value int    `"=" @( [ "-" ] Int )`
}

type Message struct {
	Pos lexer.Position

	Name    string          `"message" @Ident`
	Entries []*MessageEntry `"{" @@* "}"`
}

type MessageEntry struct {
	Pos lexer.Position

	Annotation *Annotation `@@`
	Field      *Field      `| @@  ";"*`
}

type Annotation struct {
	Pos lexer.Position

	Provider string `"+" @Ident`
	Type     string `":" @Ident ":"`

	Func  string `( @Ident`
	Value string `| @String )`
}

type Field struct {
	Pos lexer.Position

	Type *Type  `@@`
	Name string `@Ident`
	Tag  int    `"=" @Int`
}

type Scalar int

const (
	None Scalar = iota
	Double
	Float
	Int32
	Int64
	Bool
	String
	Bytes
	Datetime
)

var scalarToString = map[Scalar]string{
	None: "None", Double: "Double", Float: "Float", Int32: "Int32", Int64: "Int64",
	Bool: "Bool", String: "String", Bytes: "Bytes", Datetime: "Datetime",
}

func (s Scalar) GoString() string { return scalarToString[s] }

var stringToScalar = map[string]Scalar{
	"double": Double, "float": Float, "int32": Int32, "int64": Int64,
	"bool": Bool, "string": String, "bytes": Bytes, "datetime": Datetime,
}

func (s *Scalar) Parse(lex *lexer.PeekingLexer) error {
	token := lex.Peek()
	v, ok := stringToScalar[token.Value]
	if !ok {
		return participle.NextMatch
	}
	lex.Next()
	*s = v
	return nil
}

type Type struct {
	Pos lexer.Position

	Scalar    Scalar   `  @@`
	Map       *MapType `| @@`
	Reference string   `| @(Ident ( "." Ident )*)`
}

type MapType struct {
	Pos lexer.Position

	Key   *Type `"map" "<" @@`
	Value *Type `"," @@ ">"`
}

type ParseError struct {
	err error
}

func (p ParseError) Error() string {
	return fmt.Sprint("error at position ", p.err.(participle.Error).Position(), p.err)
}

func Parse(in []byte) (*AST, error) {
	p, err := participle.Build[Document](participle.UseLookahead(2), participle.Elide("Comment"), participle.Unquote())
	if err != nil {
		return nil, ParseError{err}
	}

	d, err := p.Parse("", bytes.NewBuffer(in))
	if err != nil {
		return nil, ParseError{err}
	}

	return toAST(*d), nil
}
