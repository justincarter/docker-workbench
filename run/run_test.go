package run

import (
	"testing"
)

func TestRun(t *testing.T) {
	err := Run("echo", "test", "pass")
	if err != nil {
		t.Fail()
	}
}

func TestOutput(t *testing.T) {
	out, err := Output("echo", "test", "pass")
	if len(out) == 0 || err != nil {
		t.Fail()
	}

}
