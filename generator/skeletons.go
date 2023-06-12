package generator

import (
	"github.com/dave/jennifer/jen"
)

func (g *Generator) generateSkeletonValidation(t string, fn string) jen.Code {
	return jen.Func().Params(jen.Id("sf").Id(t)).Id(toCamel(fn)).Params(jen.Id("string"), jen.Id("any")).Params(jen.Id("error")).Block(
		jen.Return(jen.Id("nil")),
	)
}

func (g *Generator) generateSkeletonTransform(t string, fn string) jen.Code {
	return jen.Func().Id(toCamel(fn)).Params(jen.Id("any")).Params(jen.Id("any"), jen.Id("error")).Block(
		jen.Return(jen.Id("nil"), jen.Id("nil")),
	)
}
