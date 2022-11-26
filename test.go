package sky

import (
	"math"
	"testing"
)

const precision = 0.00000001

func numberEquals(a, b interface{}) bool {
	fa, oka := a.(float64)
	fb, okb := b.(float64)
	if oka && okb {
		return math.Abs(fa-fb) < precision
	}
	return a == b
}

func assert(t *testing.T, v bool, msg string) {
	if !v {
		if msg != "" {
			t.Errorf("msg")
		}
		t.Errorf("assert error")
	}
}
