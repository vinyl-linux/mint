package parser

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	suffix = ".mint"
)

type Document struct {
	Pos lexer.Position

	Entries []*Entry `( @@ ";"* )*`
}

type Entry struct {
	Pos lexer.Position

	Type *Type ` @@`
	Enum *Enum `| @@`
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

func (e Enum) name() string {
	return e.Name
}

func (e Enum) pos() lexer.Position {
	return e.Pos
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

type Type struct {
	Pos lexer.Position

	Name    string          `"type" @Ident`
	Entries []*MessageEntry `"{" @@* "}"`
}

func (t Type) name() string {
	return t.Name
}

func (t Type) pos() lexer.Position {
	return t.Pos
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

	DataType *DataType `@@`
	Name     string    `@Ident`
	Tag      int       `"=" @Int`
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

type DataType struct {
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

var p = participle.MustBuild[Document](participle.UseLookahead(2), participle.Elide("Comment"), participle.Unquote())

func Parse(fn string, in io.Reader) (*AST, error) {
	d, err := p.Parse(fn, in)
	if err != nil {
		return nil, err
	}

	return toAST(*d)
}

func ParseDir(dir string) (*AST, error) {
	asts := make([]*AST, 0)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, suffix) {
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			defer f.Close()

			a, err := Parse(path, f)
			if err != nil {
				return err
			}

			asts = append(asts, a)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return merge(asts)
}
