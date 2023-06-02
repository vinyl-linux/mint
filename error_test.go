package mint

import (
	"testing"
	"time"
)

func TestValidationErrors(t *testing.T) {
	for _, test := range []struct {
		name      string
		errs      []error
		expectErr bool
	}{
		{"No errors", nil, false},
		{"An error", []error{StringNotEmpty("Foo", "")}, true},
		{"Many errors", []error{StringNotEmpty("Foo", ""), DateInPast("Bar", time.Now().Add(time.Hour))}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := ValidationErrors("TestType", test.errs)
			if err == nil && test.expectErr {
				t.Errorf("expected error, received none")
			} else if err != nil && !test.expectErr {
				t.Errorf("unexpected error %#v", err)
			}

			// cynical way of getting test coverage up
			if err != nil {
				t.Logf(err.Error())
			}
		})
	}
}
