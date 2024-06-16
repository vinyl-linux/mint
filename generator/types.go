package generator

import (
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vinyl-linux/mint/parser"
)

func (g *Generator) generateForType(t parser.AnnotatedType) (err error) {
	g.customFunctions = make([]jen.Code, 0)

	ret := jen.NewFile(g.PackageName)

	ret.Add(g.generateTypeDefinition(t))
	ret.Add(g.generateValidations(t))
	ret.Add(g.generateTransformations(t))
	ret.Add(g.generateValuer(t.Name))

	// We need to manually run these loops, rather than exploding
	// the output of generateUnmarshaller (etc.) as per:
	//
	//  ret.Add(g.generateUnmarshaller(t)...)
	//
	// because doing so creates invalid code for whatever
	// reason
	for _, u := range g.generateUnmarshaller(t) {
		ret.Add(u)
	}

	for _, u := range g.generateMarshaller(t) {
		ret.Add(u)
	}

	err = g.writeSkeletons(t.Name)
	if err != nil {
		return
	}

	return ret.Save(filepath.Join(g.Directory, strings.Join([]string{strings.ToLower(t.Name), "go"}, ".")))
}

// generateTypeDefinition creates the top level struct from names, types, and
// doc strings
func (g *Generator) generateTypeDefinition(at parser.AnnotatedType) (c jen.Code) {
	fields := make([]jen.Code, 0)
	for _, f := range at.Entries {
		if len(f.DocString) > 0 {
			fields = append(fields, jen.Null().Comment(f.DocString))
		}

		fields = append(fields, jen.Null().Id(f.Name).Add(toJenType(f.Field)))
	}

	return jen.Null().Type().Id(at.Name).Struct(
		fields...,
	)
}
