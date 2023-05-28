// nolint: govet, golint
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

type DataType struct {
	Pos lexer.Position

	Map    *MapType `@@`
	Scalar *Scalar  `| @@`
}

type MapType struct {
	Pos lexer.Position

	Key   string `"map" "<" @Ident`
	Value string `"," @Ident ">"`
}

type Scalar struct {
	Pos lexer.Position

	Type string `@Ident`
}

var p = participle.MustBuild[Document](participle.UseLookahead(2), participle.Elide("Comment"), participle.Unquote())

func parse(fn string, in io.Reader) (*AST, error) {
	d, err := p.Parse(fn, in)
	if err != nil {
		return nil, err
	}

	return toAST(*d)
}

func Parse(fn string, in io.Reader) (*AST, error) {
	a, err := parse(fn, in)
	if err != nil {
		return nil, err
	}

	// Because the solver does some magic around uniqueness
	// we need to call it here for a len 1 slice
	return merge([]*AST{a})
}

func ParseFile(fn string) (*AST, error) {
	// #nosec: G304
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	// #nosec: G307
	defer f.Close()

	return Parse(fn, f)
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
			// #nosec: G304
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			// #nosec: G307
			defer f.Close()

			a, err := parse(path, f)
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
