package sky

import (
	"testing"
)

func TestLiteralPrimary(t *testing.T) {
	s := `
	1 1.2 "hello" true false
`
	p := getParser(s)
	values := []interface{}{1, 1.2, "hello", true, false}
	for i := 0; i < len(values); i++ {
		v := p.primary()
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
	p := getParser(s)
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
		t.Errorf("token should be MINUS.")
	}
	literal, ok := innerValue.expression.(*literalExpr)
	if !ok || literal.value != 10 {
		t.Errorf("expect Literal '10'")
	}
}

//test association
func TestAdd(t *testing.T) {
	s := `1 + 2 - 3`
	p := getParser(s)
	expr1 := p.add()
	exprRoot, ok := expr1.(*binaryExpr)
	if !ok {
		t.Errorf("result is not binary expression.")
	}
	assert(t, exprRoot.operator.tokenType == MINUS, "left association, root should be minus")

	leftRoot, ok := exprRoot.left.(*binaryExpr)
	if !ok {
		t.Errorf("result is not binary expression.")
	}
	assert(t, leftRoot.operator.tokenType == PLUS, "left association, left root should be PLUS")
}

//test precedence
func TestAddAndMultiply(t *testing.T) {
	s := `1 * 2 - 3 % 5`
	//			-
	//       *      %
	//   1    2   3    5
	p := getParser(s)
	exp := p.add()
	exprRoot, ok := exp.(*binaryExpr)
	if !ok {
		t.Errorf("root is not binary expression.")
	}
	assert(t, exprRoot.operator.tokenType == MINUS, "root operator should be '-'")

	leftRoot, ok := exprRoot.left.(*binaryExpr)
	if !ok {
		t.Errorf("left root is not binary expression.")
	}
	assert(t, leftRoot.operator.tokenType == STAR, "left operator should be '*'")

	rightExpr, ok := exprRoot.right.(*binaryExpr)
	if !ok {
		t.Errorf("right root is not binary expression.")
	}
	assert(t, rightExpr.operator.tokenType == MOD, "right operator should be '%'")
}

func getParser(source string) *parser {
	l := NewLexer(source)
	tokens := l.parse()
	return newParser(tokens)
}
