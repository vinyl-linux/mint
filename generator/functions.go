package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

func toGoFuncName(custom bool, s string) *jen.Statement {
	csFuncName := toCamel(s)

	switch custom {
	case false:
		return jen.Qual(mintPath, csFuncName)

	default:
		return jen.Id("sf").Dot(csFuncName)

	}
}

func toCamel(s string) string {
	return strcase.ToCamel(s)
}

func marshallerFuncName(s string) string {
	return fmt.Sprintf("marshall%s", s)
}

func unmarshallerFuncName(s string) string {
	return fmt.Sprintf("unmarshall%s", s)
}
