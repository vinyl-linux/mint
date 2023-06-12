package generator

import (
	"testing"
)

func TestGenerator_generateSkeletonValidation(t *testing.T) {
	expect := `func (sf TestType) ValidateSomeField(string, any) error {
	return nil
}`

	g := new(Generator)
	received := codeToString(g.generateSkeletonValidation("TestType", "validate_some_field"))

	if expect != received {
		t.Errorf("expected\n%s\nreceived\n%s", expect, received)
	}
}

func TestGenerator_generateSkeletonTransform(t *testing.T) {
	expect := `func ChangeAThing(any) (any, error) {
	return nil, nil
}`

	g := new(Generator)
	received := codeToString(g.generateSkeletonTransform("TestType", "change_a_thing"))

	if expect != received {
		t.Errorf("expected\n%s\nreceived\n%s", expect, received)
	}
}
