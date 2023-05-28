package cmd

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}
	}()

	initConfig()
}
