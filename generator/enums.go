package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vinyl-linux/mint/parser"
)

func (g *Generator) generateForEnum(t parser.Enum) (err error) {
	ret := jen.NewFile(g.PackageName)

	// Create type as a uint8
	ret.Add(g.generateEnumDefinition(t))

	// Create values
	ret.Add(g.generateEnumValues(t))

	// Create marshaller, unmarshaller, valuer
	ret.Add(g.marshallEnum(t.Name))
	ret.Add(g.unmarshallEnum(t))
	ret.Add(g.generateValuer(t.Name))

	return ret.Save(filepath.Join(g.Directory, strings.Join([]string{strings.ToLower(t.Name), "go"}, ".")))
}

func (g *Generator) generateEnumDefinition(t parser.Enum) jen.Code {
	return jen.Null().Type().Id(t.Name).Id("byte")
}

func (g *Generator) generateEnumValues(t parser.Enum) jen.Code {
	consts := make([]jen.Code, len(t.Values)+1)

	consts[0] = jen.Id(enumValueString(t.Name, "Unknown")).Id(t.Name).Op("=").Id("iota")

	for idx, ev := range t.Values {
		consts[idx+1] = jen.Id(enumValue(t.Name, ev))
	}

	return jen.Null().Const().Defs(consts...)
}

func enumValue(en string, ee *parser.EnumEntry) string {
	return enumValueString(en, ee.Value.Key)
}

func enumValueString(en, ee string) string {
	return fmt.Sprintf("%s%s", en, ee)
}
