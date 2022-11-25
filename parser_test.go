package sky

import (
	"testing"
)

func TestLiteralPrimary(t *testing.T) {
	s := `
	1 1.2 "hello" true false
`
	l := NewLexer(s)
	tokens := l.parse()

	p := newParser(tokens)
	values := []interface{}{1, 1.2, "hello", true, false}
	for i := 0; i < len(values); i++ {
		v := p.unary()
		value, ok := v.(*literalExpr)
		if !ok {
			t.Errorf("index %d is not literalExpr", i)
		}
		if value.value != values[i] {
			t.Errorf("expected: %q, actual: %q", values[i], v)
		}
	}
}

func TestUnary(t *testing.T) {
	s := `--10`
	l := NewLexer(s)
	tokens := l.parse()

	p := newParser(tokens)
	expr := p.unary()
	value, ok := expr.(*unaryExpr)
	if !ok {
		t.Errorf("expect unaryExpr")
	}
	if value.token.tokenType != MINUS {
		t.Errorf("token should bu MINUS.")
	}

	innerValue, ok := value.expression.(*unaryExpr)
	if !ok {
		t.Errorf("expect unaryExpr")
	}
	if innerValue.token.tokenType != MINUS {
		t.Errorf("token should bu MINUS.")
	}
	literal, ok := innerValue.expression.(*literalExpr)
	if !ok || literal.value != 10 {
		t.Errorf("expect Literal '10'")
	}

}

