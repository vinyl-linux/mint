package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/vinyl-linux/mint/parser"
)

func (g Generator) marshallSliceArray(t string, e parser.AnnotatedEntry) jen.Code {
	var (
		dt string
	)

	if e.Field.DataType == nil {
		return jen.Null()
	}

	switch {
	case e.Field.DataType.Slice != nil:
		dt = e.Field.DataType.Slice.Type

	case e.Field.DataType.FixedSizeSlice != nil:
		dt = e.Field.DataType.FixedSizeSlice.Type

	default:
		return jen.Null()
	}

	fn := marshallerFuncName(e.Name)
	innerInitialiser, _, _ := scalarToMintJen(dt)

	return jen.Func().Params(jen.Id("sf").Id(t)).Id(fn).Params(jen.Id("w").Qual("io", "Writer")).Params(jen.Id("err").Id("error")).
		Block(
			jen.Id("f").Op(":=").Id("make").Call(jen.Index().Add(muvType), jen.Id("len").Call(jen.Id("sf").Dot(e.Field.Name))),
			jen.For(jen.Id("i").Op(":=").Id("range").Id("f")).Block(
				jen.Id("f").Index(jen.Id("i")).Op("=").Add(innerInitialiser).Call(jen.Id("sf").Dot(e.Field.Name).Index(jen.Id("i"))),
			),
			jen.Return(jen.Qual(mintPath, "NewSliceCollection").Call(jen.Id("f"), jen.Lit(e.Field.DataType.FixedSizeSlice != nil)).Dot("Marshall").Call(jen.Id("w"))),
		)

}

func (g Generator) marshallMap(t string, e parser.AnnotatedEntry) jen.Code {
	fn := marshallerFuncName(e.Name)

	keyInitialiser, _, _ := scalarToMintJen(e.DataType.Map.Key)
	valueInitialiser, _, _ := scalarToMintJen(e.DataType.Map.Value)

	return jen.Func().Params(jen.Id("sf").Id(t)).Id(fn).Params(jen.Id("w").Qual("io", "Writer")).Params(jen.Id("err").Id("error")).
		Block(
			jen.Id("f").Op(":=").Id("make").Call(jen.Id("map").Index(jen.Add(muvType)).Add(muvType)),
			jen.For(jen.Id("k").Op(",").Id("v").Op(":=").Id("range").Id("sf").Dot(e.Name)).Block(
				jen.Id("f").Index(jen.Add(keyInitialiser).Call(jen.Id("k"))).Op("=").Add(valueInitialiser).Call(jen.Id("v")),
			),
			jen.Return(jen.Qual(mintPath, "NewMapCollection").Call(jen.Id("f")).Dot("Marshall").Call(jen.Id("w"))),
		)

}
