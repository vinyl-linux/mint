package generator

import (
	"os"
	"testing"

	"github.com/vinyl-linux/mint/parser"
	"golang.org/x/mod/sumdb/dirhash"
)

func TestNew(t *testing.T) {
	opts := &GeneratorOptions{
		PackageName: "testtest",
	}

	g, _ := New(nil, opts)

	if g.PackageName != "testtest" {
		t.Errorf("expected testtest, received %s", g.PackageName)
	}
}

// TestGenerator_Generate is a really silly test, really;
//
// We essentially generate a directory full of code and then
// compare some checksums- this test _wont_ tell you where things
// have failed; you need other tests to fail to tell you that.
//
// This test serves two purposes:
//
//  1. Act as a canary for the efficacy of tests (if this fails and
//     nothing else fails then we have unexpected behaviour)
//  2. Allows us to provide test coverage for the largest function;
//     the main generator
func TestGenerator_Generate(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}

	ast, err := parser.ParseDir("testdata/valid-documents")
	if err != nil {
		t.Fatal(err)
	}

	opts := &GeneratorOptions{
		PackageName:             "mint_testrun",
		MakeDirectory:           true,
		Directory:               dir,
		CustomFunctionSkeletons: true,
		Clobber:                 true,
	}

	g, _ := New(ast, opts)
	err = g.Generate()
	if err != nil {
		t.Fatal(err)
	}

	expect := "h1:H/ivkG97ylTC8lCmAZbtmzstj+/foFqaNEKOiLRIdcQ="
	received, err := dirhash.HashDir(dir, "", dirhash.DefaultHash)
	if err != nil {
		t.Fatal(err)
	}

	if expect != received {
		t.Fatalf("expected %q, received %q", expect, received)
	}
}

func TestGenerator_generateType(t *testing.T) {
	g := new(Generator)

	expect := `type SomeTestType struct {
	SomeStringSlice []string
	SomeStringSlice [5]string
	// An int64 for some reason
	SomeStringSlice map[string]int64
	// ATypeOfSomeType is a uuid
	ATypeOfSomeType v5.UUID
	Thingy          BlahType
}`
	received := codeToString(g.generateType(simpleType))

	if expect != received {
		t.Errorf("expected\n%s\nreceived\n%s", expect, received)
	}

}

func TestGenerator_generateValidations(t *testing.T) {
	g := new(Generator)
	g.GeneratorOptions.CustomFunctionSkeletons = true

	expectFuncs := `package test

func (sf SomeTestType) BlahBlahBlahHowDoesThisEvaluate(string, any) error {
	return nil
}
`

	expect := `func (sf SomeTestType) Validate() error {
	errors := make([]error, 0)
	for _, err := range []error{mint.NotEmpty("ATypeOfSomeType", sf.ATypeOfSomeType), sf.BlahBlahBlahHowDoesThisEvaluate("ATypeOfSomeType", sf.ATypeOfSomeType)} {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return mint.ValidationErrors("SomeTestType", errors)
}`
	t.Run("Validate()", func(t *testing.T) {
		received := codeToString(g.generateValidations(simpleType))

		if expect != received {
			t.Errorf("expected\n%s\nreceived\n%s", expect, received)
		}
	})

	t.Run("Skeleton function(s)", func(t *testing.T) {
		received := codeSliceToFile(g.customFunctions)

		if expectFuncs != received {
			t.Errorf("expected\n%s\nreceived\n%s", expectFuncs, received)
		}

	})
}

func TestGenerator_generateTransformations(t *testing.T) {
	g := new(Generator)
	g.GeneratorOptions.CustomFunctionSkeletons = true

	expectFuncs := `package test

func Floop(any) (any, error) {
	return nil, nil
}
func TrebleValue(any) (any, error) {
	return nil, nil
}
`

	expect := `func (sf *SomeTestType) Transform() (err error) {
	sf.SomeStringSlice, err = sf.Floop(sf.SomeStringSlice)
	if err != nil {
		return
	}
	sf.SomeStringSlice, err = sf.TrebleValue(sf.SomeStringSlice)
	if err != nil {
		return
	}
	sf.SomeStringSlice, err = mint.Flipbits(sf.SomeStringSlice)
	if err != nil {
		return
	}
	return
}`
	t.Run("Validate()", func(t *testing.T) {
		received := codeToString(g.generateTransformations(simpleType))

		if expect != received {
			t.Errorf("expected\n%s\nreceived\n%s", expect, received)
		}
	})

	t.Run("Skeleton function(s)", func(t *testing.T) {
		received := codeSliceToFile(g.customFunctions)

		if expectFuncs != received {
			t.Errorf("expected\n%s\nreceived\n%s", expectFuncs, received)
		}

	})
}

func TestGenerator_generateUnmarshaller(t *testing.T) {
	g := new(Generator)

	expect := simpleTypeStringUnmarshaller
	received := codeSliceToFile(g.generateUnmarshaller(simpleType))

	if expect != received {
		t.Errorf("expected\n%s\nreceived\n%s", expect, received)
	}
}

func TestGenerator_generateMarshaller(t *testing.T) {
	g := new(Generator)

	expect := simpleTypeStringMarshaller
	received := codeSliceToFile(g.generateMarshaller(simpleType))

	if expect != received {
		t.Errorf("expected\n%s\nreceived\n%s", expect, received)
	}
}

func TestGenerator_generateValuer(t *testing.T) {
	g := new(Generator)

	expect := `func (sf SomeTestType) Value() any {
	return sf
}`
	received := codeToString(g.generateValuer(simpleType))

	if expect != received {
		t.Errorf("expected\n%s\nreceived\n%s", expect, received)
	}
}

func TestScalarToMintJen(t *testing.T) {
	for _, test := range []struct {
		ts                string
		expectInitialiser string
		expectNilValue    string
		expectCastType    string
	}{
		{"string", "mint.NewStringScalar", `""`, "string"},
		{"datetime", "mint.NewDatetimeScalar", "time.Time{}", "time.Time"},
		{"uuid", "mint.NewUuidScalar", "v5.UUID{}", "v5.UUID"},
		{"uint32", "mint.NewUInt32Scalar", "uint32(0)", "uint32"},
		{"int16", "mint.NewInt16Scalar", "int16(0)", "int16"},
		{"int32", "mint.NewInt32Scalar", "int32(0)", "int32"},
		{"int64", "mint.NewInt64Scalar", "int64(0)", "int64"},
		{"float32", "mint.NewFloat32Scalar", "float32(0)", "float32"},
		{"float64", "mint.NewFloat64Scalar", "float64(0)", "float64"},
		{"bool", "mint.NewBoolScalar", "false", "bool"},
		{"byte", "mint.NewByteScalar", "byte(int32(0))", "byte"},
		{"SomeType", "new", "SomeType", "SomeType"},
	} {
		t.Run(test.ts, func(t *testing.T) {
			receivedInitialiser, receivedNilValue, receivedCastType := scalarToMintJen(test.ts)

			t.Run("initialiser", func(t *testing.T) {
				s := codeToString(receivedInitialiser)
				if test.expectInitialiser != s {
					t.Errorf("expected %s, reveived %s", test.expectInitialiser, s)
				}
			})

			t.Run("nilvalue", func(t *testing.T) {
				s := codeToString(receivedNilValue)
				if test.expectNilValue != s {
					t.Errorf("expected %s, reveived %s", test.expectNilValue, s)
				}
			})

			t.Run("castType", func(t *testing.T) {
				s := codeToString(receivedCastType)
				if test.expectCastType != s {
					t.Errorf("expected %s, reveived %s", test.expectCastType, s)
				}
			})

		})
	}
}
