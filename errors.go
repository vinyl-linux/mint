package mint

import (
	"fmt"
	"strings"
)

type ErrValidationErrors struct {
	t    string
	errs []error
}

func (e ErrValidationErrors) Error() string {
	out := strings.Builder{}
	out.WriteString(fmt.Sprintf("instance of type %s failed validation:\n", e.t))

	for _, err := range e.errs {
		out.WriteString("\t")
		out.WriteString(err.Error())
		out.WriteString("\n")
	}

	return out.String()
}

func ValidationErrors(t string, errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	return ErrValidationErrors{
		t:    t,
		errs: errs,
	}
}
