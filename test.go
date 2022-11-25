package sky

import "testing"

func assert(t *testing.T, v bool, msg string) {
	if !v {
		if msg != "" {
			t.Errorf("msg")
		}
		t.Errorf("assert error")
	}
}
