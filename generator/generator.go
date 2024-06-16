package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vinyl-linux/mint/parser"
)

const (
	mintPath = "github.com/vinyl-linux/mint"
)

var (
	muvType = jen.Qual(mintPath, "MarshallerUnmarshallerValuer")
)

type GeneratorOptions struct {
	PackageName             string
	MakeDirectory           bool
	Directory               string
	CustomFunctionSkeletons bool
	Clobber                 bool
}

type Generator struct {
	GeneratorOptions

	ast             *parser.AST
	customFunctions []jen.Code
}

func New(doc *parser.AST, options *GeneratorOptions) (g *Generator, err error) {
	g = new(Generator)
	g.GeneratorOptions = *options
	g.ast = doc

	return
}

// Generate will create:
//
//  1. Type definitions
//  2. Enums
//  3. Validations
//  4. Transforms
//  5. Unmarshallers; and
//  6. Marshallers
//
// For each type defined in an AST
func (g *Generator) Generate() (err error) {
	if g.MakeDirectory {
		err = os.MkdirAll(g.Directory, 0750)
		if err != nil {
			return
		}
	}

	for _, t := range g.ast.Types {
		err = g.generateForType(t)
		if err != nil {
			return
		}
	}

	for _, e := range g.ast.Enums {
		err = g.generateForEnum(e)
		if err != nil {
			return
		}
	}

	return
}

func (g Generator) writeSkeletons(tn string) (err error) {
	if !g.CustomFunctionSkeletons {
		return
	}

	fn := filepath.Join(g.Directory, strings.Join([]string{strings.ToLower(tn), "custom", "go"}, "."))
	customs := jen.NewFile(g.PackageName)
	for _, f := range g.customFunctions {
		customs.Add(f)
	}

	if _, statErr := os.Stat(fn); statErr != nil || g.Clobber {
		err = customs.Save(fn)
		if err != nil {
			return
		}
	}

	return
}

// generateValidations creates calls to both mint and custom
// validations, additionally templating custom validations were
// requested
func (g *Generator) generateValidations(at parser.AnnotatedType) (c jen.Code) {
	functionCalls := make([]jen.Code, 0)
	for _, e := range at.Entries {
		for _, f := range e.Validations {
			fn := toGoFuncName(f.IsCustom, f.Function)

			if g.CustomFunctionSkeletons && f.IsCustom {
				g.customFunctions = append(g.customFunctions, g.generateSkeletonValidation(at.Name, f.Function))
			}

			functionCalls = append(functionCalls, fn.Call(jen.Lit(e.Name), jen.Id("sf").Dot(e.Name)))
		}
	}

	return jen.Func().Params(jen.Id("sf").Id(at.Name)).
		Id("Validate").Params().Params(jen.Id("error")).
		Block(
			jen.Id("errors").Op(":=").Id("make").Call(jen.Index().Id("error"), jen.Lit(0)),
			jen.For(jen.List(jen.Id("_"), jen.Id("err")).Op(":=").Range().Index().Id("error").Values(
				functionCalls...,
			)).
				Block(
					jen.If(
						jen.Id("err").Op("!=").Id("nil")).Block(
						jen.Id("errors").Op("=").Id("append").Call(
							jen.Id("errors"), jen.Id("err")),
					),
				),
			jen.Return().Id("mint").Dot("ValidationErrors").Call(jen.Lit(at.Name), jen.Id("errors")),
		)
}

// generateTransformations creates calls to both mint and custom
// transformations, additionally templating custom transformations were
// requested
func (g *Generator) generateTransformations(at parser.AnnotatedType) jen.Code {
	functionCalls := make([]jen.Code, 0)
	for _, e := range at.Entries {
		for _, f := range e.Transformations {
			fn := toGoFuncName(f.IsCustom, f.Function)

			if g.CustomFunctionSkeletons && f.IsCustom {
				g.customFunctions = append(g.customFunctions, g.generateSkeletonTransform(at.Name, f.Function))
			}

			functionCalls = append(functionCalls,
				jen.List(jen.Id("sf").Dot(e.Name), jen.Id("err")).Op("=").Add(fn).Call(jen.Id("sf").Dot(e.Name)),
				jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
					jen.Return(),
				),
			)
		}
	}

	functionCalls = append(functionCalls, jen.Return())

	return jen.Func().Params(jen.Id("sf").Op("*").Id(at.Name)).Id("Transform").Params().Params(jen.Id("err").Id("error")).
		Block(
			functionCalls...,
		)
}

// generateUnmarshaller will:
//  1. Create an unmarshall function per field
//  2. Create an implementation of the mint.Unmarshaller interface for this type
func (g *Generator) generateUnmarshaller(at parser.AnnotatedType) (j []jen.Code) {
	functionCalls := make([]jen.Code, 0)
	j = make([]jen.Code, 0)

	for _, entry := range at.Entries {
		switch {
		case entry.DataType.Scalar != nil:
			j = append(j, g.unmarshallScalar(at.Name, entry))

		case entry.DataType.Slice != nil,
			entry.DataType.FixedSizeSlice != nil:
			j = append(j, g.unmarshallSliceArray(at.Name, entry))

		case entry.DataType.Map != nil:
			j = append(j, g.unmarshallMap(at.Name, entry))

		default:
			continue
		}

		fn := jen.Id("sf").Dot(unmarshallerFuncName(entry.Name))
		functionCalls = append(functionCalls,
			jen.If(jen.Id("err").Op("=").Add(fn).Call(jen.Id("r")).Id(";").Id("err").Op("!=").Id("nil")).Block(
				jen.Return(),
			),
		)
	}

	functionCalls = append(functionCalls,
		callErrorable("Transform"),
		callErrorable("Validate"),
		jen.Return(),
	)

	j = append(j, jen.Func().Params(jen.Id("sf").Op("*").Id(at.Name)).Id("Unmarshall").Params(jen.Id("r").Qual("io", "Reader")).Params(jen.Id("err").Id("error")).
		Block(
			functionCalls...,
		))

	return
}

// generateMarshaller will:
//  1. Create a marshall function per field
//  2. Create an implementation of the mint.Marshaller interface for this type
func (g *Generator) generateMarshaller(at parser.AnnotatedType) (j []jen.Code) {
	functionCalls := []jen.Code{
		callErrorable("Transform"),
		callErrorable("Validate"),
	}

	j = make([]jen.Code, 0)

	for _, e := range at.Entries {
		fieldName := e.Field.Name

		switch {
		// if a scalar, and creator isn't 'new' then do
		// err = mint.$creator(value).Marshall(w) (etc)
		// else just value.Marshall(w)
		case e.DataType.Scalar != nil:
			if _, ok := parser.Scalars[e.DataType.Scalar.Type]; ok {
				f, _, _ := scalarToMintJen(e.DataType.Scalar.Type)
				functionCalls = append(functionCalls,
					jen.If(jen.Id("err").Op("=").Add(f).Call(jen.Id("sf").Dot(fieldName)).Dot("Marshall").Call(jen.Id("w")).Id(";").Id("err").Op("!=").Id("nil")).Block(jen.Return()),
				)

				continue
			}

			functionCalls = append(functionCalls,
				jen.If(jen.Id("err").Op("=").Id("sf").Dot(fieldName).Dot("Marshall").Call(jen.Id("w")).Id(";").Id("err").Op("!=").Id("nil")).Block(jen.Return()),
			)

		// if a slice or array, then mint.NewSlicecollection, setting
		// the second arg accordingly)
		case e.DataType.Slice != nil ||
			e.DataType.FixedSizeSlice != nil:
			j = append(j, g.marshallSliceArray(at.Name, e))

			fn := jen.Id("sf").Dot(marshallerFuncName(e.Name))
			functionCalls = append(functionCalls,
				jen.If(jen.Id("err").Op("=").Add(fn).Call(jen.Id("w")).Id(";").Id("err").Op("!=").Id("nil")).Block(jen.Return()),
			)

		// if a map then go from one to the other
		case e.DataType.Map != nil:
			j = append(j, g.marshallMap(at.Name, e))

			fn := jen.Id("sf").Dot(marshallerFuncName(e.Name))
			functionCalls = append(functionCalls,
				jen.If(jen.Id("err").Op("=").Add(fn).Call(jen.Id("w")).Id(";").Id("err").Op("!=").Id("nil")).Block(jen.Return()),
			)

		default:
			continue
		}
	}

	functionCalls = append(functionCalls, jen.Return())

	return append(j, jen.Func().Params(jen.Id("sf").Id(at.Name)).Id("Marshall").Params(jen.Id("w").Qual("io", "Writer")).Params(jen.Id("err").Id("error")).
		Block(
			functionCalls...,
		),
	)
}

// generateValuer will create an implementation of the mint.Valuer interface
// for this type.
//
// This valuer is, more or less, a no-op; it's there to implement the interface
// mint.MarshallerUnmarshallerValuer that allows us to work with binary
// representations of more complex, or exciting types
func (g *Generator) generateValuer(n string) jen.Code {
	return jen.Func().Params(jen.Id("sf").Id(n)).Id("Value").Params().Params(jen.Id("any")).
		Block(
			jen.Return(jen.Id("sf")),
		)
}

func toJenType(t parser.Field) jen.Code {
	dt := t.DataType

	switch {
	case dt.Scalar != nil:
		_, _, goType := scalarToMintJen(t.DataType.Scalar.Type)

		return goType

	case dt.Slice != nil:
		return jen.Index().Id(dt.Slice.Type)

	case dt.FixedSizeSlice != nil:
		return jen.Index(jen.Id(fmt.Sprintf("%d", dt.FixedSizeSlice.Size))).Id(dt.FixedSizeSlice.Type)

	case dt.Map != nil:
		return jen.Map(jen.Id(dt.Map.Key)).Id(dt.Map.Value)
	}

	return jen.Null()
}

func callErrorable(f string) jen.Code {
	return jen.If(jen.Id("err").Op("=").Id("sf").Dot(f).Call().Id(";").Id("err").Op("!=").Id("nil")).Block(jen.Return())
}

func scalarToMintJen(ts string) (c, nilValue, castType jen.Code) {
	switch ts {
	case "string":
		return jen.Qual(mintPath, "NewStringScalar"), jen.Lit(""), jen.Id(ts)

	case "datetime":
		return jen.Qual(mintPath, "NewDatetimeScalar"), jen.Qual("time", "Time").Block(), jen.Qual("time", "Time")

	case "uuid":
		return jen.Qual(mintPath, "NewUuidScalar"), jen.Qual("github.com/gofrs/uuid/v5", "UUID").Block(), jen.Qual("github.com/gofrs/uuid/v5", "UUID")

	case "uint32":
		return jen.Qual(mintPath, "NewUInt32Scalar"), jen.Id("uint32").Call(jen.Lit(0)), jen.Id(ts)

	case "int16":
		return jen.Qual(mintPath, "NewInt16Scalar"), jen.Id("int16").Call(jen.Lit(0)), jen.Id(ts)

	case "int32":
		return jen.Qual(mintPath, "NewInt32Scalar"), jen.Id("int32").Call(jen.Lit(0)), jen.Id(ts)

	case "int64":
		return jen.Qual(mintPath, "NewInt64Scalar"), jen.Id("int64").Call(jen.Lit(0)), jen.Id(ts)

	case "float32":
		return jen.Qual(mintPath, "NewFloat32Scalar"), jen.Id("float32").Call(jen.Lit(0)), jen.Id(ts)

	case "float64":
		return jen.Qual(mintPath, "NewFloat64Scalar"), jen.Id("float64").Call(jen.Lit(0)), jen.Id(ts)

	case "bool":
		return jen.Qual(mintPath, "NewBoolScalar"), jen.Lit(false), jen.Id(ts)

	case "byte", "uint8":
		return jen.Qual(mintPath, "NewByteScalar"), jen.Id("byte").Call(jen.Lit('\x00')), jen.Id(ts)
	}

	return jen.Id("new"), jen.Id(ts), jen.Id(ts)
}
