package mint

import (
	"testing"
	"time"
)

func TestDateInPast(t *testing.T) {
	for _, test := range []struct {
		name      string
		t         time.Time
		expectErr bool
	}{
		{"Future date", time.Now().Add(time.Hour * 24 * 365), true},
		{"Past date", time.Now().Add(0 - time.Hour), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := DateInPast(test.name, test.t)
			if err == nil && test.expectErr {
				t.Errorf("expected error, received none")
			} else if err != nil && !test.expectErr {
				t.Errorf("unexpected error %#v", err)
			}
		})
	}
}

func TestDateInUTC(t *testing.T) {
	in := time.Now().In(time.FixedZone("Seoul", 9*60*60))
	if s, _ := in.Zone(); s != "Seoul" {
		t.Fatalf("expected Local, received %s", s)
	}

	out, err := DateInUtc(in)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	if s, _ := out.Zone(); s != "UTC" {
		t.Fatalf("expected Local, received %s", s)
	}
}

func TestStringNotEmpty(t *testing.T) {
	for _, test := range []struct {
		s         string
		expectErr bool
	}{
		{"", true},
		{"hello, world!", false},
	} {
		t.Run(test.s, func(t *testing.T) {
			err := StringNotEmpty(test.s, test.s)
			if err == nil && test.expectErr {
				t.Errorf("expected error, received none")
			} else if err != nil && !test.expectErr {
				t.Errorf("unexpected error %#v", err)
			}
		})
	}
}
