package parser

import (
	"fmt"
	"sort"

	"github.com/alecthomas/participle/v2/lexer"
)

// missingTagErr describes occasions where a type has a
// missing tag; all tags (even deprecated ones) must still
// exist in a message or things break
type missingTagErr struct {
	t   string
	tag int
}

// Error returns an error message describing which tag is
// missing in which type
func (e missingTagErr) Error() string {
	return fmt.Sprintf("missing tag %d for type %s",
		e.tag,
		e.t,
	)
}

// merge takes a slice of asts, ensures uniqueness of names, and
// returns either an error describing collisions, or the union of
// all ASTs
func merge(in []*AST) (out *AST, err error) {
	names := make(map[string][]lexer.Position)
	intermediateOut := new(AST)

	for _, a := range in {
		if a == nil {
			continue
		}

		intermediateOut.Types = append(intermediateOut.Types, a.Types...)
		intermediateOut.Enums = append(intermediateOut.Enums, a.Enums...)

		for _, t := range namedSlice(a.Types, a.Enums) {
			if _, ok := names[t.name()]; !ok {
				names[t.name()] = make([]lexer.Position, 0)
			}

			names[t.name()] = append(names[t.name()], t.pos())
		}
	}

	err = toCollisionError("type", names)
	if err != nil {
		return
	}

	// Validate all type fields against types/ enums, scalars map
	for _, t := range intermediateOut.Types {
		for _, e := range t.Entries {
			err = e.IsValidType(names)
			if err != nil {
				return
			}
		}
	}

	return intermediateOut, nil
}

// toAnnotatedType accepts a Type definiton from our parser, and
// generates an annotated type.
//
// To wit:
//  1. Iterate through the entries in a type
//  2. Determine which validations, transforms, and doc strings belong to which 'thing'
//  3. Ensure each entry is unique in name
//  4. Ensure each entry has a unique position tag
//
// Groupings occur by parsing each entry until we hit a field definition, and then
// merging those entries into a single definition.
func toAnnotatedType(m Type) (a AnnotatedType, err error) {
	a.Pos = m.Pos
	a.Name = m.Name
	a.Entries = make([]AnnotatedEntry, 0)

	names := make(map[string][]lexer.Position)
	tags := make(map[string][]lexer.Position)
	tagValues := make([]int, 0)

	ae := AnnotatedEntry{}
	for _, e := range m.Entries {
		if e.Annotation != nil {
			switch e.Annotation.Type {
			case "doc":
				ae.AppendDocString(e.Annotation.Value)
			case "validate":
				ae.AppendValidation(Validation{
					IsCustom: e.Annotation.Provider == "custom",
					Function: e.Annotation.Func,
				})
			case "transform":
				ae.AppendTransformation(Transformation{
					IsCustom: e.Annotation.Provider == "custom",
					Function: e.Annotation.Func,
				})
			}

			continue
		}

		// If we get here, we get to a field and so are finishing this
		// annotated message
		if _, ok := names[e.Field.Name]; !ok {
			names[e.Field.Name] = make([]lexer.Position, 0)
		}

		names[e.Field.Name] = append(names[e.Field.Name], e.Field.Pos)

		tagStr := intToStr(e.Field.Tag)
		if _, ok := tags[tagStr]; !ok {
			tags[tagStr] = make([]lexer.Position, 0)
		}

		tags[tagStr] = append(tags[tagStr], e.Field.Pos)

		tagValues = append(tagValues, e.Field.Tag)

		ae.Field = *e.Field
		a.Entries = append(a.Entries, ae)
		ae = AnnotatedEntry{}
	}

	// Ensure names are unique
	err = toCollisionError(a.Name+" field", names)
	if err != nil {
		return
	}

	// Ensure tags are unique
	err = toCollisionError(a.Name+" tag value", tags)
	if err != nil {
		return
	}

	// Ensure there are no missing tags
	is := sort.IntSlice(tagValues)
	is.Sort()

	for idx, v := range is {
		if idx != v {
			err = missingTagErr{
				t:   a.Name,
				tag: idx,
			}

			return
		}
	}

	return
}

func intToStr(i int) string {
	return fmt.Sprintf("%d", i)
}
