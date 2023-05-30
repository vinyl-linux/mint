package mint

import (
	"fmt"
	"time"
)

func DateInPast(name string, t time.Time) error {
	now := time.Now().UTC().UnixNano()
	if now < t.UTC().UnixNano() {
		return fmt.Errorf("%s should be in the past", t)
	}
	return nil
}

func DateInUTC(t time.Time) (time.Time, error) {
	return t.UTC(), nil
}

func StringNotEmpty(name, s string) error {
	if len(s) == 0 {
		return fmt.Errorf("%s should not be empty", name)
	}
	return nil
}
