package parser

type AST struct {
	Messages []AnnotatedMessage
	Enums    []Enum
}

func toAST(d Document) (ad *AST) {
	ad = new(AST)
	ad.Enums = make([]Enum, 0)
	ad.Messages = make([]annotatedMessage, 0)

	for _, e := range d.Entries {
		if e.Enum != nil {
			ad.Enums = append(ad.Enums, *e.Enum)

			continue
		}

		ad.Messages = append(ad.Messages, ToAnnotatedMessage(*e.Message))
	}

	return
}

type annotatedMessage struct {
	Name    string
	Entries []annotatedEntry
}

type annotatedEntry struct {
	Field

	DocString       string
	Validations     []Validation
	Transformations []Transformation
}

func (ae *annotatedEntry) AppendDocString(s string) {
	if len(ae.DocString) == 0 {
		ae.DocString = s

		return
	}
	ae.DocString = ae.DocString + " " + s
}

func (ae *annotatedEntry) AppendValidation(v Validation) {
	if ae.Validations == nil {
		ae.Validations = make([]Validation, 0)
	}

	ae.Validations = append(ae.Validations, v)
}

func (ae *annotatedEntry) AppendTransformation(v Transformation) {
	if ae.Transformations == nil {
		ae.Transformations = make([]Transformation, 0)
	}

	ae.Transformations = append(ae.Transformations, v)
}

type validation struct {
	IsCustom bool
	Function string
}

type transformation struct {
	IsCustom bool
	Function string
}

func toAnnotatedMessage(m Message) (a annotatedMessage) {
	a.Name = m.Name
	a.Entries = make([]annotatedEntry, 0)

	ae := annotatedEntry{}
	for _, e := range m.Entries {
		if e.Annotation != nil {
			switch e.Annotation.Type {
			case "doc":
				ae.AppendDocString(e.Annotation.Value)
			case "validate":
				ae.AppendValidation(validation{
					IsCustom: e.Annotation.Provider == "custom",
					Function: e.Annotation.Func,
				})
			case "transform":
				ae.AppendTransformation(transformation{
					IsCustom: e.Annotation.Provider == "custom",
					Function: e.Annotation.Func,
				})
			}

			continue
		}

		// If we get here, we get to a field and so are finishing this
		// annotated message
		ae.Field = *e.Field
		a.Entries = append(a.Entries, ae)
		ae = annotatedEntry{}
	}

	return
}
