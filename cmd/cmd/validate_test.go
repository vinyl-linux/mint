package cmd

import (
	"testing"
)

type returnedStatus struct {
	v int
}

func (r *returnedStatus) fail(i int) { r.v = i }

func TestValidateCmd_Run(t *testing.T) {
	origFailer := failer
	defer func() {
		failer = origFailer
	}()

	for _, test := range []struct {
		dir          string
		expectStatus int
	}{
		{"testdata/valid-documents", 0},
		{"testdata/invalid-docuemnt", 1},
		{"testdata/nonsuch", 1},
	} {
		t.Run(test.dir, func(t *testing.T) {
			rs := returnedStatus{0}
			failer = rs.fail
			validateCmd.Run(nil, []string{test.dir})

			if test.expectStatus != rs.v {
				t.Errorf("expected %d, received %d", test.expectStatus, rs.v)
			}
		})
	}
}

func TestValidateCmd_Arg(t *testing.T) {
	for _, test := range []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{"missing arg errors", []string{}, true},
		{"two many args errors", []string{"dir1/", "dir2/"}, true},
		{"one arg is successful", []string{"dir1/"}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := validateCmd.Args(nil, test.args)
			if err == nil && test.expectErr {
				t.Errorf("expected error, received none")
			} else if err != nil && !test.expectErr {
				t.Errorf("unexpected error %#v", err)
			}
		})
	}
}
