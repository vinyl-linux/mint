package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/vinyl-linux/mint/parser"
)

func (g Generator) unmarshallSliceArray(t string, e parser.AnnotatedEntry) jen.Code {
	fn := unmarshallerFuncName(e.Name)
	var (
		block []jen.Code
		maker jen.Code
		dt    string
	)

	if e.Field.DataType == nil {
		return jen.Null()
	}

	switch {
	case e.Field.DataType.Slice != nil:
		block = unmarshallSlicePreludeGetLen(e)
		dt = e.Field.DataType.Slice.Type
		maker = jen.Id("sf").Dot(e.Field.Name).Op("=").Id("make").Call(jen.Index().Id(dt), jen.Id("f").Dot("Len").Call())

	case e.Field.DataType.FixedSizeSlice != nil:
		block = unmarshallSlicePreludeFixedLen(e)
		dt = e.Field.DataType.FixedSizeSlice.Type
		maker = jen.Id("sf").Dot(e.Field.Name).Op("=").Index(jen.Lit(e.Field.DataType.FixedSizeSlice.Size)).Id(dt).Block()

	default:
		return jen.Null()
	}

	innerInitialiser, innerNilValue, innerCastType := scalarToMintJen(dt)

	block = append(block, []jen.Code{
		jen.For(jen.List(jen.Id("i"), jen.Null()).Op(":=").Range().Id("f").Dot("V")).Block(
			jen.Id("f").Dot("V").Index(jen.Id("i")).Op("=").Add(innerInitialiser).Call(innerNilValue),
		),
		jen.Id("err").Op("=").Id("f").Dot("Unmarshall").Call(jen.Id("r")),
		jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
			jen.Return(),
		),
		maker,
		jen.For(
			jen.List(jen.Id("i"), jen.Id("v")).Op(":=").Range().Id("f").Dot("Value").Call().Assert(jen.Index().Qual(mintPath, "MarshallerUnmarshallerValuer"))).Block(
			jen.Id("sf").Dot(e.Field.Name).Index(jen.Id("i")).Op("=").Id("v").Dot("Value").Call().Assert(innerCastType),
		),
		jen.Return(),
	}...)

	return jen.Func().Params(jen.Id("sf").Op("*").Id(t)).Id(fn).Params(jen.Id("r").Qual("io", "Reader")).Params(jen.Id("err").Id("error")).
		Block(
			block...,
		)
}

func (g Generator) unmarshallMap(t string, e parser.AnnotatedEntry) jen.Code {
	fn := unmarshallerFuncName(e.Name)

	keyInitialiser, keyNilValue, keyCastType := scalarToMintJen(e.DataType.Map.Key)
	valueInitialiser, valueNilValue, valueCastType := scalarToMintJen(e.DataType.Map.Value)

	return jen.Func().Params(jen.Id("sf").Op("*").Id(t)).Id(fn).Params(jen.Id("r").Qual("io", "Reader")).Params(jen.Id("err").Id("error")).
		Block(
			jen.Id("f").Op(":=").Qual(mintPath, "NewMapCollection").Call(jen.Id("map").Index(muvType).Add(muvType).Block()),
			jen.Id("err").Op("=").Id("f").Dot("ReadSize").Call(jen.Id("r")),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return(),
			),

			jen.For(jen.Id("i").Op(":=").Lit(0), jen.Id("i").Op("<").Id("f").Dot("Len").Call(), jen.Id("i").Op("++")).Block(
				jen.Id("f").Dot("V").Index(jen.Add(keyInitialiser).Call(keyNilValue)).Op("=").Add(valueInitialiser).Call(valueNilValue),
			),

			jen.Id("err").Op("=").Id("f").Dot("Unmarshall").Call(jen.Id("r")),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return(),
			),

			jen.Id("sf").Dot(e.Name).Op("=").Id("make").Call(jen.Map(keyCastType).Add(valueCastType)),

			jen.For(jen.List(jen.Id("k"), jen.Id("v")).Op(":=").Range().Id("f").Dot("Value").Call().Assert(jen.Map(muvType).Add(muvType))).Block(
				jen.Id("sf").Dot(e.Name).Index(jen.Id("k").Dot("Value").Call().Assert(keyCastType)).Op("=").Id("v").Dot("Value").Call().Assert(valueCastType),
			),
			jen.Return(),
		)

}

func (g Generator) unmarshallScalar(t string, e parser.AnnotatedEntry) jen.Code {
	fn := unmarshallerFuncName(e.Name)
	initialiser, nilValue, castType := scalarToMintJen(e.Field.DataType.Scalar.Type)

	return jen.Func().Params(jen.Id("sf").Op("*").Id(t)).Id(fn).Params(jen.Id("r").Qual("io", "Reader")).Params(jen.Id("err").Id("error")).
		Block(
			jen.Id("f").Op(":=").Add(initialiser).Call(nilValue),
			jen.Id("err").Op("=").Id("f").Dot("Unmarshall").Call(jen.Id("r")),
			jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
				jen.Return(),
			),
			jen.Id("sf").Dot(e.Name).Op("=").Id("f").Dot("Value").Call().Assert(castType),
			jen.Return(),
		)
}

func unmarshallSlicePreludeGetLen(e parser.AnnotatedEntry) []jen.Code {
	return []jen.Code{
		jen.Id("f").Op(":=").Id("mint").Dot("NewSliceCollection").Call(jen.Id("nil"), jen.Id("false")),
		jen.Id("err").Op("=").Id("f").Dot("ReadSize").Call(jen.Id("r")),
		jen.If(jen.Id("err").Op("!=").Id("nil")).Block(
			jen.Return(),
		),
		jen.Id("f").Dot("V").Op("=").Id("make").Call(jen.Index().Qual(mintPath, "MarshallerUnmarshallerValuer"), jen.Id("f").Dot("Len").Call()),
	}
}

func unmarshallSlicePreludeFixedLen(e parser.AnnotatedEntry) []jen.Code {
	return []jen.Code{jen.Id("f").Op(":=").Qual(mintPath, "NewSliceCollection").Call(jen.Id("make").Call(jen.Index().Add(muvType), jen.Lit(e.DataType.FixedSizeSlice.Size)), jen.Id("true"))}
}
